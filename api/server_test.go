package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerRejectsInvalidAPIRoutes(t *testing.T) {
	server := newServer()
	for _, test := range []struct {
		name   string
		method string
		path   string
		status int
	}{
		{"wrong method", http.MethodGet, "/api/game/id/deal", http.StatusMethodNotAllowed},
		{"unknown path", http.MethodGet, "/api/unknown", http.StatusNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			server.handler().ServeHTTP(response, httptest.NewRequest(test.method, test.path, nil))
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
			if strings.Contains(response.Body.String(), "<!DOCTYPE html>") {
				t.Fatal("API request returned SPA HTML")
			}
		})
	}
}
