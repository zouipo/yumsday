package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

// IntPathValues is a middleware that parses integer values from the URL path
// and stores them in the request context.
func IntPathValues(valueNames ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = parsePathValue(
				w,
				r,
				func(valueName, valueStr string) (any, error) {
					value, err := strconv.ParseInt(valueStr, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("%s must be a valid interger", valueName)
					}
					return value, nil
				},
				valueNames...,
			)
			if r == nil {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// FloatPathValues is a middleware that parses float values from the URL path
// and stores them in the request context.
func FloatPathValues(valueNames ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = parsePathValue(
				w,
				r,
				func(valueName, valueStr string) (any, error) {
					value, err := strconv.ParseFloat(valueStr, 64)
					if err != nil {
						return nil, fmt.Errorf("%s must be a valid floating point number", valueName)
					}
					return value, nil
				},
				valueNames...,
			)
			if r == nil {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// StringPathValues is a middleware that parses string values from the URL path
// and stores them in the request context.
func StringPathValues(valueNames ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = parsePathValue(
				w,
				r,
				func(valueName, valueStr string) (any, error) {
					return valueStr, nil
				},
				valueNames...,
			)
			if r == nil {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

/*** PRIVATE HELPERS ***/

// pathValueParser defines a function type for parsing path values from strings to their expected types.
type pathValueParser func(valueName, valueStr string) (any, error)

// parsePathValue extracts and parses values from the URL path of an HTTP request.
// It uses the provided parseFunc to convert string values to their expected types (int, float, string).
// Parsed values are stored in the request context for later use by handlers.
func parsePathValue(
	w http.ResponseWriter,
	r *http.Request,
	parseFunc pathValueParser,
	valueNames ...string,
) *http.Request {
	// Loop through every value name to parse from the URL path.
	for _, valueName := range valueNames {
		valueStr := r.PathValue(valueName)
		if valueStr == "" {
			http.Error(
				w,
				fmt.Sprintf("Failed to parse %s from request URL", valueName),
				http.StatusBadRequest,
			)
			return nil
		}

		// Parse the value from string to its expected type, using the provided parse function.
		value, err := parseFunc(valueName, valueStr)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return nil
		}

		// Store the parsed value in the request context for final handler use.
		r = r.WithContext(context.WithValue(r.Context(), valueName, value))
		slog.Debug(
			fmt.Sprintf("Parsed %s from URL", valueName),
			valueName,
			value,
		)
	}

	return r
}
