package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/service"
)

func SessionInjector(sessionService service.SessionServiceInterface) Middleware {
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
				Name:  sessionService.CookieName(),
				Value: s.ID,
				// JS cannot access the cookie via document.cookie;
				// security measure that prevents XSS attacks from stealing the session ID
				HttpOnly: true,
				// The URL path prefix for which the browser will send the cookie.
				// "/" means it is sent on all paths of the domain.
				Path: "/",
				// The browser will only send the cookie over HTTPS connections;
				// prevents the session ID from being intercepted over plain HTTP.
				Secure: true,
				// The cookie will be sent by the browser only when Same-Site Sub-Requests are done,
				// or when the domain is typed directly in the URL.
				// It won't be sent for cross-site sub-requests, which helps mitigate CSRF attacks.
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(sessionService.Expiration()).UTC(),
				// The lifetime of the cookie in seconds from when it was received.
				// Prefered over Expires because it is not dependent on the client's clock,
				// but we set them both for compatibility with older browsers.
				MaxAge: int(sessionService.Expiration().Seconds()),
			}
			// Adds a Set-Cookie header to the ResponseWriter's headers.
			// This header instructs the browser to store the cookie and its attributes
			// and send it with future requests to the same domain.
			http.SetCookie(w, cookie)

			next.ServeHTTP(w, r)

			if r.URL.Path != "/logout" {
				// Save session in dedicated goroutine to reduce response latency.
				go sessionService.Save(s)
			}
		})
	}
}
