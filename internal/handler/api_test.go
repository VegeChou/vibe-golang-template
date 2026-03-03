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

func decodeAPIResponse(t *testing.T, body *bytes.Buffer) response.APIResponse[any] {
	t.Helper()

	var out response.APIResponse[any]
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

func TestCreateUserUnknownFieldRejected(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/users",
		bytes.NewReader([]byte(`{"name":"Alice","email":"alice@example.com","role":"admin"}`)),
	)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	out := decodeAPIResponse(t, resp.Body)
	if out.Code != response.CodeCommonInvalidParam {
		t.Fatalf("expected COMMON_INVALID_PARAM, got %s", out.Code)
	}
}

func TestCreateAndListUsersWithPageResult(t *testing.T) {
	mux := newTestMux()

	for _, user := range []map[string]string{
		{"name": "Alice", "email": "alice@example.com"},
		{"name": "Bob", "email": "bob@example.com"},
	} {
		body, _ := json.Marshal(user)
		createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
		createResp := httptest.NewRecorder()
		mux.ServeHTTP(createResp, createReq)
		if createResp.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d", createResp.Code)
		}
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&size=1", nil)
	listResp := httptest.NewRecorder()
	mux.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.Code)
	}

	listOut := decodeAPIResponse(t, listResp.Body)
	if !listOut.Success || listOut.Code != response.CodeOK {
		t.Fatalf("expected success OK envelope, got success=%v code=%s", listOut.Success, listOut.Code)
	}

	dataMap, ok := listOut.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be object")
	}
	if int(dataMap["page"].(float64)) != 1 || int(dataMap["size"].(float64)) != 1 {
		t.Fatalf("expected page=1,size=1, got page=%v,size=%v", dataMap["page"], dataMap["size"])
	}
	if int(dataMap["total"].(float64)) != 2 {
		t.Fatalf("expected total=2, got %v", dataMap["total"])
	}
}

func TestListUsersInvalidPageSize(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=0&size=101", nil)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	out := decodeAPIResponse(t, resp.Body)
	if out.Code != response.CodeCommonInvalidParam {
		t.Fatalf("expected COMMON_INVALID_PARAM, got %s", out.Code)
	}
}
