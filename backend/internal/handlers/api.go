package handlers

import (
	"database/sql"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zouipo/yumsday/backend/internal/middleware"
	"github.com/zouipo/yumsday/backend/internal/repositories"
	"github.com/zouipo/yumsday/backend/internal/services"
	// _ "github.com/zouipo/yumsday/backend/docs"
)

// NewAPIServer registers API routes on a new ServeMux.
func NewAPIServer(db *sql.DB) http.Handler {
	// ServeMux = HTTP request multiplexer, a router.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern tha most closely matches the URL.
	mux := http.NewServeMux()

	// Swagger = provides a UI for API documentation
	mux.Handle("/swagger/", httpSwagger.Handler())

	// Initializing every layers
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)
	userHandler.RegisterRoutes(mux, "/api/user")

	middlewareStack := middleware.Stack(
		middleware.ResponseWritter,
		middleware.Logger,
	)

	return middlewareStack(mux)
}
