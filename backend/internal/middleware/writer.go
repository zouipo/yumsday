package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

// Custom response writer to intercept calls to WriterHeader
// and Writter to do additional processing, e.g get the status code.
type ResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// WriteHeader intercepts the call to WriteHeader to capture the status code.
func (w *ResponseWriter) WriteHeader(status int) {
	// The wroteHeader flag is to ensure we only call WriteHeader once,
	// otherwise we get a "superfluous response.WriteHeader call" error in the logs.
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.status = status
	w.ResponseWriter.WriteHeader(status)
	slog.Debug("WriteHeader", "status", status)
}

// Write intercepts the call to Write (for example by json.Encode) to ensure WriteHeader is called first.
// http.Error explicitely calls our WriteHeader (check http.Error) sources,
// so it's not bypassed in that case.
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	return w.ResponseWriter.Write(data)
}

// ResponseWritter is a middleware that wraps the ResponseWriter struct to capture status codes.
func ResponseWritter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		// Store a pointer to the status so the logger can read the updated value
		r = r.WithContext(context.WithValue(r.Context(), "status", &writer.status))
		next.ServeHTTP(writer, r)
	})
}
