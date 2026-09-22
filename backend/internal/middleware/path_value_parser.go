package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/zouipo/yumsday/backend/internal/ctxkey"
)

// IntPathValues is a middleware that parses integer values from the URL path
// and stores them in the request context.
func IntPathValues(keys ...ctxkey.Key) Middleware {
	return newParserMiddleWare(intParser, keys...)
}

// FloatPathValues is a middleware that parses float values from the URL path
// and stores them in the request context.
func FloatPathValues(keys ...ctxkey.Key) Middleware {
	return newParserMiddleWare(floatParser, keys...)
}

// StringPathValues is a middleware that parses string values from the URL path
// and stores them in the request context.
func StringPathValues(keys ...ctxkey.Key) Middleware {
	return newParserMiddleWare(stringParser, keys...)
}

func IdPathValue() Middleware {
	return IntPathValues(ctxkey.Id{})
}

/*** PRIVATE HELPERS ***/

// pathValueParser defines a function type for parsing path values from strings to their expected types.
type pathValueParser func(key ctxkey.Key, valueStr string) (any, error)

func intParser(key ctxkey.Key, valueStr string) (any, error) {
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s must be a valid integer", key)
	}
	return value, nil
}

func floatParser(key ctxkey.Key, valueStr string) (any, error) {
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, fmt.Errorf("%s must be a valid floating point number", key)
	}
	return value, nil
}

func stringParser(key ctxkey.Key, valueStr string) (any, error) {
	return valueStr, nil
}

func newParserMiddleWare(parser pathValueParser, keys ...ctxkey.Key) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = parsePathValue(w, r, parser, keys...)
			if r == nil {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// parsePathValue extracts and parses values from the URL path of an HTTP request.
// It uses the provided parseFunc to convert string values to their expected types (int, float, string).
// Parsed values are stored in the request context for later use by handler.
func parsePathValue(
	w http.ResponseWriter,
	r *http.Request,
	parseFunc pathValueParser,
	keys ...ctxkey.Key,
) *http.Request {
	for _, key := range keys {
		valueStr := r.PathValue(key.String())
		if valueStr == "" {
			http.Error(
				w,
				fmt.Sprintf("Failed to parse %s from request URL", key),
				http.StatusBadRequest,
			)
			return nil
		}

		value, err := parseFunc(key, valueStr)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return nil
		}

		r = r.WithContext(context.WithValue(r.Context(), key, value))
		slog.Debug(
			fmt.Sprintf("Parsed %s from URL", key),
			key,
			value,
		)
	}

	return r
}
