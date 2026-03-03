package response

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type ErrorCode string

const (
	CodeOK                   ErrorCode = "OK"
	CodeCommonInvalidParam   ErrorCode = "COMMON_INVALID_PARAM"
	CodeAuthUnauthorized     ErrorCode = "AUTH_UNAUTHORIZED"
	CodeAuthForbidden        ErrorCode = "AUTH_FORBIDDEN"
	CodeCommonNotFound       ErrorCode = "COMMON_NOT_FOUND"
	CodeCommonConflict       ErrorCode = "COMMON_CONFLICT"
	CodeCommonTooManyRequest ErrorCode = "COMMON_TOO_MANY_REQUESTS"
	CodeDependencyError      ErrorCode = "SYSTEM_DEPENDENCY_ERROR"
	CodeInternalError        ErrorCode = "SYSTEM_INTERNAL_ERROR"

	DefaultLang = "en-US"
	LangZhCN    = "zh-CN"
	LangEnUS    = "en-US"
)

type APIResponse struct {
	Success   bool      `json:"success"`
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Lang      string    `json:"lang"`
	Data      any       `json:"data"`
	TraceID   string    `json:"traceId"`
	Timestamp string    `json:"timestamp"`
}

func Message(lang, en, zh string) string {
	if lang == LangZhCN {
		return zh
	}
	return en
}

func ResolveLang(r *http.Request) string {
	acceptLang := strings.TrimSpace(r.Header.Get("Accept-Language"))
	if strings.HasPrefix(acceptLang, LangZhCN) {
		return LangZhCN
	}
	if strings.HasPrefix(acceptLang, LangEnUS) {
		return LangEnUS
	}

	queryLang := strings.TrimSpace(r.URL.Query().Get("lang"))
	if queryLang == LangZhCN || queryLang == LangEnUS {
		return queryLang
	}

	return DefaultLang
}

func Write(w http.ResponseWriter, r *http.Request, status int, success bool, code ErrorCode, message string, data any) {
	body := APIResponse{
		Success:   success,
		Code:      code,
		Message:   message,
		Lang:      ResolveLang(r),
		Data:      data,
		TraceID:   traceIDFromRequest(r),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, r *http.Request, status int, data any, message string) {
	Write(w, r, status, true, CodeOK, message, data)
}

func Error(w http.ResponseWriter, r *http.Request, status int, code ErrorCode, message string) {
	Write(w, r, status, false, code, message, nil)
}

func traceIDFromRequest(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("X-Trace-Id")); v != "" {
		return v
	}

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "trace-id-unavailable"
	}

	return hex.EncodeToString(buf)
}
