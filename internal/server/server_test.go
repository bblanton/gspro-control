package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func setupTestServer() *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}

func TestHealth(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rw.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestActions(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/actions", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var body map[string][]string
	if err := json.Unmarshal(rw.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := body["mulligan"]; !ok {
		t.Fatalf("expected default action 'mulligan' present; got %v", reflect.ValueOf(body).MapKeys())
	}
}

func TestCommand_MissingAction(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/command", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestCommand_UnknownAction(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/command?action=notreal", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
}

func TestCommand_Success(t *testing.T) {
	prev := executeComboFunc
	executeComboFunc = func(combo []string) error { return nil }
	defer func() { executeComboFunc = prev }()

	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/command?action=mulligan", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var body commandResponse
	if err := json.Unmarshal(rw.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Result != "success" || body.Action != "mulligan" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestCommand_Error(t *testing.T) {
	prev := executeComboFunc
	executeComboFunc = func(combo []string) error { return errors.New("boom") }
	defer func() { executeComboFunc = prev }()

	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/command?action=mulligan", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rw.Code)
	}
	var body commandResponse
	if err := json.Unmarshal(rw.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Result != "error" || body.Action != "mulligan" || body.Error == "" {
		t.Fatalf("unexpected body: %#v", body)
	}
}
