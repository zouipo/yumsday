package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/ctxkey"
)

type expectedStatusKey struct{}

// final handler
func writerReqHandler(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value(expectedStatusKey{}).(int)
	w.WriteHeader(status)
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

	handler.ServeHTTP(w, r)

	if statusExtractor.status != expectedStatus {
		t.Errorf("expected status %v, got %v", expectedStatus, statusExtractor.status)
	}
}

func TestWriter(t *testing.T) {
	statusExtractor := &statusExtractor{}
	mwStack := Stack(ResponseWriter, statusExtractor.middleware)
	handler := mwStack(http.HandlerFunc(writerReqHandler))

	performRequestAndCheck(t, statusExtractor, handler, http.StatusOK)
	performRequestAndCheck(t, statusExtractor, handler, http.StatusInternalServerError)
}
