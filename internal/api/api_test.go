package api

import (
	"example.com/noise-registry/internal/service"
	"example.com/noise-registry/internal/store"
	"net/http/httptest"
	"testing"
)

func TestHTTPHealth(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	New(service.New(s)).Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
}
