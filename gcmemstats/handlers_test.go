package gcmemstats

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	header = "text/plain; version=0.0.4; charset=utf-8"
)

func TestMetricsHandler(t *testing.T) {
	handler := MetricsHandler()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, header, w.Header().Get("Content-Type"))

	body := w.Body.String()
	assert.Contains(t, body, "gc_alloc_bytes")
	assert.Contains(t, body, "gc_num_gc")
	assert.Contains(t, body, "gc_num_goroutine")
	assert.Contains(t, body, "# HELP")
	assert.Contains(t, body, "# TYPE")
}

func TestSetGCPercent(t *testing.T) {
	originalPercent := SetGCPercent(100)

	newPercent := SetGCPercent(150)
	assert.Equal(t, 100, newPercent)

	SetGCPercent(originalPercent)
}

func TestRegisterPprof(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPprof(mux)

	endpoints := []string{
		"/debug/pprof/",
		"/debug/pprof/heap",
		"/debug/pprof/goroutine",
		"/debug/pprof/allocs",
	}

	for _, endpoint := range endpoints {
		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNotFound, w.Code, "Endpoint %s should be registered", endpoint)
	}
}
