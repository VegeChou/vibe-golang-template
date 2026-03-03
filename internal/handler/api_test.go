package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vibe-golang-template/internal/repository/memory"
	"vibe-golang-template/internal/service"
	"vibe-golang-template/pkg/response"
)

func newTestMux() *http.ServeMux {
	repo := memory.NewUserRepository()
	svc := service.NewUserService(repo)
	api := NewAPI(svc)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	return mux
}

func decodeAPIResponse(t *testing.T, body *bytes.Buffer) response.APIResponse {
	t.Helper()

	var out response.APIResponse
	if err := json.Unmarshal(body.Bytes(), &out); err != nil {
		t.Fatalf("failed to decode api response: %v", err)
	}
	return out
}

func TestHealthzEnvelope(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz?lang=zh-CN", nil)
	req.Header.Set("X-Trace-Id", "trace-1")
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	body := decodeAPIResponse(t, resp.Body)
	if !body.Success {
		t.Fatalf("expected success true")
	}
	if body.Code != response.CodeOK {
		t.Fatalf("expected code OK, got %s", body.Code)
	}
	if body.Lang != response.LangZhCN {
		t.Fatalf("expected lang zh-CN, got %s", body.Lang)
	}
	if body.TraceID != "trace-1" {
		t.Fatalf("expected trace id trace-1, got %s", body.TraceID)
	}
	if body.Timestamp == "" {
		t.Fatalf("expected timestamp not empty")
	}
}

func TestCreateUserInvalidInput(t *testing.T) {
	mux := newTestMux()
	body, _ := json.Marshal(map[string]string{"name": "", "email": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Accept-Language", "en-US")
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	out := decodeAPIResponse(t, resp.Body)
	if out.Success {
		t.Fatalf("expected success false")
	}
	if out.Code != response.CodeCommonInvalidParam {
		t.Fatalf("expected COMMON_INVALID_PARAM, got %s", out.Code)
	}
	if out.Data != nil {
		t.Fatalf("expected null data for error response")
	}
}

func TestCreateAndListUsers(t *testing.T) {
	mux := newTestMux()

	body, _ := json.Marshal(map[string]string{"name": "Alice", "email": "alice@example.com"})
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	createReq.Header.Set("Accept-Language", "en-US")
	createResp := httptest.NewRecorder()
	mux.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResp.Code)
	}

	createOut := decodeAPIResponse(t, createResp.Body)
	if !createOut.Success || createOut.Code != response.CodeOK {
		t.Fatalf("expected success OK envelope, got success=%v code=%s", createOut.Success, createOut.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	listReq.Header.Set("Accept-Language", "en-US")
	listResp := httptest.NewRecorder()
	mux.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.Code)
	}

	listOut := decodeAPIResponse(t, listResp.Body)
	if !listOut.Success || listOut.Code != response.CodeOK {
		t.Fatalf("expected success OK envelope, got success=%v code=%s", listOut.Success, listOut.Code)
	}
}
