package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vibe-golang-template/internal/repository/memory"
	"vibe-golang-template/internal/service"
)

func newTestMux() *http.ServeMux {
	repo := memory.NewUserRepository()
	svc := service.NewUserService(repo)
	api := NewAPI(svc)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	return mux
}

func TestHealthz(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestCreateAndListUsers(t *testing.T) {
	mux := newTestMux()

	body, _ := json.Marshal(map[string]string{"name": "Alice", "email": "alice@example.com"})
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	createResp := httptest.NewRecorder()
	mux.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResp.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	listResp := httptest.NewRecorder()
	mux.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.Code)
	}
}
