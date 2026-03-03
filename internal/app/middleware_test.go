package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vibe-golang-template/pkg/response"
)

func TestRecoverMiddlewareReturnsUnifiedError(t *testing.T) {
	h := recoverMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.Code)
	}

	var out response.APIResponse[any]
	if err := json.Unmarshal(resp.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if out.Success {
		t.Fatalf("expected success=false")
	}
	if out.Code != response.CodeInternalError {
		t.Fatalf("expected SYSTEM_INTERNAL_ERROR, got %s", out.Code)
	}
	if out.Data != nil {
		t.Fatalf("expected null data")
	}
}
