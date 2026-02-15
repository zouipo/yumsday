package middleware

import (
	"net/http"
	"slices"
)

// A middleware in Go is fundamentally a function that wraps an http.Handler to
// perform additional processing on HTTP requests and responses.
// It takes and returns a http.Handler, allowing chaining of multiple middlewares.
type Middleware func(http.Handler) http.Handler

// Stack combines multiple middlewares into a single middleware
// that executes them in the order they are passed.
//
// Middleware execution follows an onion-like pattern where each middleware wraps
// the next one. The execution flow is:
//
//	Request Flow (inward):
//	  Client → MW1 → MW2 → MW3 → Final Handler
//
//	Response Flow (outward):
//	  Client ← MW1 ← MW2 ← MW3 ← Final Handler
//
// Parameters:
//   - middlewares: variable number of middleware functions to be chained together
//
// Returns:
//   - A single Middleware that applies all provided middlewares in order
func Stack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// Iterate backwards through the middlewares to build the nested chain.
		// This ensures that the first middleware in the list becomes the outermost
		// layer (executed first on request, last on response), and the last
		// middleware becomes the innermost layer (executed last on request, first
		// on response).
		//
		// Example with MW1, MW2, MW3:
		//   Iteration 1: next = MW3(finalHandler)
		//   Iteration 2: next = MW2(MW3(finalHandler))
		//   Iteration 3: next = MW1(MW2(MW3(finalHandler)))
		for _, m := range slices.Backward(middlewares) {
			next = m(next)
		}
		return next
	}
}
