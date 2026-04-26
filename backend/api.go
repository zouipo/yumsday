package backend

import (
	"database/sql"
	"io/fs"
	"net/http"
	"sync"
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
func NewAPIServer(db *sql.DB, migrationsFs, front fs.FS, tasksWG *sync.WaitGroup) http.Handler {
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
	sessionInjector := middleware.SessionInjector(sessionService, tasksWG)

	authService := service.NewAuthService(sessionService, userService)
	authHandler := handler.NewAuthHandler(authService)

	middlewareStack := middleware.Stack(
		middleware.ResponseWriter,
		middleware.Logger,
		sessionInjector,
		middleware.UserInjector(userService),
	)

	swaggerMiddlewareStack := middleware.Stack(
		middleware.ResponseWriter,
		middleware.Logger,
	)

	// ServeMux = HTTP request multiplexer, a router.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern tha most closely matches the URL.
	mux := http.NewServeMux()
	backMux := http.NewServeMux()

	mux.Handle("/", http.FileServerFS(front))
	mux.Handle("/swagger/", swaggerMiddlewareStack(httpSwagger.Handler()))
	mux.Handle("/api/", middlewareStack(backMux))
	mux.Handle("/auth/", middlewareStack(backMux))

	userHandler.RegisterRoutes(backMux, "/api/user")
	authHandler.RegisterRoutes(backMux, "/auth")

	return mux
}
