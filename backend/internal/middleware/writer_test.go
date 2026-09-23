package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/ctxkey"
)

type expectedStatusKey struct{}
type testKey struct{}

// final handler
func writerReqHandler(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value(expectedStatusKey{}).(int)
	w.WriteHeader(status)

	// write the status code as a string in the response body
	statusStr := strconv.FormatInt(int64(status), 10)
	io.WriteString(w, statusStr)
}

// middleware object used to extract the http status code
// from the request context
type statusExtractor struct {
	status int
}

func (s *statusExtractor) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// call final handler so that WriteHeader is called
		next.ServeHTTP(w, r)

		// extract status code
		s.status = *r.Context().Value(ctxkey.Status{}).(*int)
	})
}

func performRequestAndCheck(t *testing.T, statusExtractor *statusExtractor, handler http.Handler, expectedStatus int) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r = r.WithContext(context.WithValue(r.Context(), expectedStatusKey{}, expectedStatus))
	r = r.WithContext(context.WithValue(r.Context(), testKey{}, t))

	handler.ServeHTTP(w, r)

	if statusExtractor.status != expectedStatus {
		t.Errorf("expected status %v, got %v", expectedStatus, statusExtractor.status)
	}

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected no error while reading response body, got %s", err)
	}

	expectedBody := strconv.FormatInt(int64(statusExtractor.status), 10)
	if string(body) != expectedBody {
		t.Errorf("expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestWriter(t *testing.T) {
	statusExtractor := &statusExtractor{}
	mwStack := Stack(ResponseWriter, statusExtractor.middleware)
	handler := mwStack(http.HandlerFunc(writerReqHandler))

	performRequestAndCheck(t, statusExtractor, handler, http.StatusOK)
	performRequestAndCheck(t, statusExtractor, handler, http.StatusInternalServerError)
}
