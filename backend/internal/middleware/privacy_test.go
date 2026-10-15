package middleware

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerOmitsQueryValues(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	handler := Logger(zap.New(core))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "https://example.invalid/api/names?token=synthetic-sensitive-value", nil))
	entries := logs.All()
	if len(entries) != 1 {
		t.Fatal("missing request event")
	}
	fields := entries[0].ContextMap()
	if _, ok := fields["query"]; ok {
		t.Fatal("raw query recorded")
	}
	if fields["path"] != "/api/names" {
		t.Fatal("request path missing")
	}
}
