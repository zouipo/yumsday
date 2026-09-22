package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type logHandler struct {
	logged bool
	record slog.Record
}

func newLogHandler() *logHandler {
	return &logHandler{
		logged: false,
	}
}

func (l *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (l *logHandler) Handle(ctx context.Context, r slog.Record) error {
	l.logged = true
	l.record = r
	return nil
}

func (l *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return l
}

func (l *logHandler) WithGroup(name string) slog.Handler {
	return l
}

func logReqHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestLogger(t *testing.T) {
	logHandler := newLogHandler()
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	mwStack := Stack(ResponseWriter, Logger)
	handler := mwStack(http.HandlerFunc(logReqHandler))

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !logHandler.logged {
		t.Error("expected logger middleware to be called")
	}
}
