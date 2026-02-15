package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testkey is a custom type for context keys to avoid collisions.
type testkey int

// ContextValueKey is the key used to store and retrieve values from the request context in tests (random).
// This key will be shared by the middlewares and the final handler to retrieve the same data.
const (
	ContextValueKey testkey = iota
)

// TestMiddlewareStack tests the Stack function by creating a middleware stack
// with two middlewares that modify the request context. It verifies that the
// final handler receives the expected context value after all middlewares have
// processed the request.
func TestMiddlewareStack(t *testing.T) {
	middlewareContextValue := "Hello world!"

	// Final handler that checks the context value set by the middlewares.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("in handler")
		// Retrieve the context value set by the middlewares.
		value := r.Context().Value(ContextValueKey).(string)

		// Verify that the context value matches the expected value.
		if value != middlewareContextValue {
			t.Errorf("expected context value %s, got %s", middlewareContextValue, value)
		}

		w.Write([]byte(value))
	})

	// First to be called on request, last on response.
	// Sets the initial context value to "Hello".
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("in middleware 1")
			// Add "Hello" to the context, at the ContextValueKey.
			r = r.WithContext(
				context.WithValue(
					r.Context(),
					ContextValueKey,
					"Hello",
				),
			)
			next.ServeHTTP(w, r)
		})
	}

	// Second to be called on request, first on response.
	// Appends " world" to the context value.
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("in middleware2")
			// Append " world" to the existing context value at ContextValueKey.
			r = r.WithContext(
				context.WithValue(
					r.Context(),
					ContextValueKey,
					r.Context().Value(ContextValueKey).(string)+" world!",
				),
			)
			next.ServeHTTP(w, r)
		})
	}

	stack := Stack(
		middleware1,
		middleware2,
	)

	// Create a test HTTP request and response recorder (their content doesn't really matter).
	r := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	// Serve the HTTP request through the middleware stack and final handler.
	stack(handler).ServeHTTP(w, r)
}
