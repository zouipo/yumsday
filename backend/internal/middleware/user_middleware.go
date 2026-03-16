package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/service"
)

func UserInjector(userService service.UserServiceInterface) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := r.Context().Value(ctx.SessionCtxKey{}).(*model.Session)

			// Not authenticated
			if s.UserID == 0 {
				slog.Debug("session is not authenticated", "id", s.ID)
				if r.URL.Path != "/login" {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			slog.Debug("session is authenticated", "id", s.ID, "user", s.UserID)

			// Authenticated but requesting /login => redirect to root
			if r.URL.Path == "/login" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			user, err := userService.GetByID(s.UserID)
			if err != nil {
				panic(err)
			}

			slog.Debug("found session user", "id", user.ID, "username", user.Username)

			r = r.WithContext(context.WithValue(
				r.Context(),
				ctx.UserCtxKey{},
				user,
			))

			next.ServeHTTP(w, r)
		})
	}
}
