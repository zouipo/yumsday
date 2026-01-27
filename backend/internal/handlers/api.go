package handlers

import (
	"database/sql"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	// _ "github.com/zouipo/yumsday/backend/docs"
)

func NewAPIServer(db *sql.DB) {
	// ServeMux = HTTP request multiplexer, a router.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern tha most closely matches the URL.
	mux := http.NewServeMux()

	// Swagger = provides a UI for API documentation
	mux.Handle("/swagger/", httpSwagger.Handler())
}
