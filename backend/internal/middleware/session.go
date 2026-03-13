package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/service"
)

func SessionInjector(sessionService *service.SessionService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := sessionService.GetSession(r)
			// http.Request context is immutable, so we need to create a new context with the session and assign it back to the request.
			r = r.WithContext(context.WithValue(
				r.Context(),
				ctx.SessionCtxKey{},
				s,
			))

			cookie := &http.Cookie{
				Name:     sessionService.CookieName(),
				Value:    s.ID,
				Domain:   "localhost",
				HttpOnly: true,
				Path:     "/",
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(sessionService.Expiration()).UTC(),
				MaxAge:   int(sessionService.Expiration().Seconds()),
			}
			// Adds a Set-Cookie header to the ResponseWriter's headers.
			http.SetCookie(w, cookie)

			next.ServeHTTP(w, r)

			if r.URL.Path != "/logout" {
				// Save session in dedicated goroutine to reduce response latency.
				go sessionService.Save(s)
			}
		})
	}
}
