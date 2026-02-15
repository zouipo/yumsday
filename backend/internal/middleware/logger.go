package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is a middleware that logs incoming requests and their completion.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing: log the incoming request
		slog.Info("Received HTTP request", "method", r.Method, "path", r.URL.Path)

		// Record the start time to measure duration when the request processing completes.
		start := time.Now().UTC()

		// Call the next handler in the chain (could be another middleware or the final handler)
		next.ServeHTTP(w, r)

		// Post-processing: after each nested middleware/handler has processed the request,
		// the logger middleware get the result back and logs the completion.
		status := r.Context().Value("status").(*int)
		slog.Info(
			"Processed HTTP request",
			"status", *status, // captured by the custom ResponseWriter middleware
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
