package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anechytailenko/Microservices_practice04/internal/api"
)

func TestHealthCheckEndpoint(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	mux := api.NewRouter(nil, nil, nil)

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	expectedBody := `{"status":"UP"}` + "\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expectedBody)
	}
}
