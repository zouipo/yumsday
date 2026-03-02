package backend

import (
	"database/sql"
	"io/fs"
	"log/slog"
	"net/http"

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

	// ServeMux = HTTP request multiplexer, a router.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern tha most closely matches the URL.
	mux := http.NewServeMux()

	// Swagger = provides a UI for API documentation
	mux.Handle("/swagger/", httpSwagger.Handler())

	// Initializing every layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userHandler.RegisterRoutes(mux, "/api/user")

	sessionRepo := repository.NewSessionRepository(db)
	slog.Info("Session repository initialized %s", sessionRepo)

	middlewareStack := middleware.Stack(
		middleware.ResponseWritter,
		middleware.Logger,
	)

	return middlewareStack(mux)
}
