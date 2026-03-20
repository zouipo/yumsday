package backend

import (
	"database/sql"
	"io/fs"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zouipo/yumsday/backend/internal/handler"
	"github.com/zouipo/yumsday/backend/internal/middleware"
	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/repository"
	"github.com/zouipo/yumsday/backend/internal/service"
	_ "github.com/zouipo/yumsday/docs"
)

// NewAPIServer registers API routes on a new ServeMux.
func NewAPIServer(db *sql.DB, migrationsFs fs.FS) http.Handler {
	err := migration.Migrate(db, migrationsFs)
	if err != nil {
		panic(err)
	}

	// Initializing every layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	sessionRepo := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(
		sessionRepo,
		"yumsday_session",
		30*24*time.Hour,
	)

	authService := service.NewAuthService(sessionService, userService)
	authHandler := handler.NewAuthHandler(authService)

	middlewareStack := middleware.Stack(
		middleware.ResponseWritter,
		middleware.Logger,
		middleware.SessionInjector(sessionService),
		middleware.UserInjector(userService),
	)

	authMiddlewareStack := middleware.Stack(
		middleware.ResponseWritter,
		middleware.Logger,
		middleware.SessionInjector(sessionService),
	)

	swaggerMiddlewareStack := middleware.Stack(
		middleware.ResponseWritter,
		middleware.Logger,
	)

	// ServeMux = HTTP request multiplexer, a router.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern tha most closely matches the URL.
	mux := http.NewServeMux()
	apiMux := http.NewServeMux()

	// Swagger = provides a UI for API documentation
	mux.Handle("/swagger/", swaggerMiddlewareStack(httpSwagger.Handler()))
	mux.Handle("/api/", middlewareStack(apiMux))
	mux.Handle("/login", authMiddlewareStack(apiMux))
	mux.Handle("/logout", middlewareStack(apiMux))

	userHandler.RegisterRoutes(apiMux, "/api/user")
	authHandler.RegisterRoutes(apiMux)

	return mux
}
