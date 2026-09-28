package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORS_OptionsRequest(t *testing.T) {
	dummyCalled := false
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dummyCalled = true
	})

	handler := withCORS(dummyHandler)

	req := httptest.NewRequest(http.MethodOptions, "/api/currencies/today", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if dummyCalled {
		t.Error("OPTIONS request should return before calling inner handler")
	}

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("expected Access-Control-Allow-Origin 'http://localhost:5173', got %q", origin)
	}

	methods := rec.Header().Get("Access-Control-Allow-Methods")
	if methods != "GET, OPTIONS" {
		t.Errorf("expected Access-Control-Allow-Methods 'GET, OPTIONS', got %q", methods)
	}
}

func TestWithCORS_GetRequest(t *testing.T) {
	dummyCalled := false
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dummyCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pass"))
	})

	handler := withCORS(dummyHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !dummyCalled {
		t.Error("expected inner handler to be called on GET request")
	}

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("expected Access-Control-Allow-Origin 'http://localhost:5173', got %q", origin)
	}
}
