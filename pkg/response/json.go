package response

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
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

type APIResponse[T any] struct {
	Success   bool      `json:"success"`
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Lang      string    `json:"lang"`
	Data      T         `json:"data"`
	TraceID   string    `json:"traceId"`
	Timestamp string    `json:"timestamp"`
}

type Translator interface {
	Translate(lang, key string) string
}

type mapTranslator struct {
	messages map[string]map[string]string
}

func (m *mapTranslator) Translate(lang, key string) string {
	if langMap, ok := m.messages[lang]; ok {
		if message, ok := langMap[key]; ok {
			return message
		}
	}
	if defaultMap, ok := m.messages[DefaultLang]; ok {
		if message, ok := defaultMap[key]; ok {
			return message
		}
	}
	return key
}

var (
	translatorMu sync.RWMutex
	translator   Translator = &mapTranslator{messages: map[string]map[string]string{
		LangEnUS: {
			"common.success":             "success",
			"error.method_not_allowed":   "method not allowed",
			"error.invalid_json_body":    "invalid json body",
			"error.json_single_object":   "json body must contain a single object",
			"error.page_invalid":         "page must be an integer >= 1",
			"error.size_invalid":         "size must be an integer between 1 and 100",
			"error.limit_invalid":        "limit must be an integer between 1 and 100",
			"error.cursor_page_conflict": "cursor and page cannot be used together",
			"error.invalid_user_input":   "invalid user input",
			"error.internal":             "internal error",
		},
		LangZhCN: {
			"common.success":             "成功",
			"error.method_not_allowed":   "请求方法不允许",
			"error.invalid_json_body":    "JSON 请求体无效",
			"error.json_single_object":   "JSON 请求体只能包含一个对象",
			"error.page_invalid":         "page 必须是大于等于 1 的整数",
			"error.size_invalid":         "size 必须是 1 到 100 之间的整数",
			"error.limit_invalid":        "limit 必须是 1 到 100 之间的整数",
			"error.cursor_page_conflict": "cursor 与 page 不能同时使用",
			"error.invalid_user_input":   "用户输入不合法",
			"error.internal":             "系统内部错误",
		},
	}}
)

type APIError struct {
	Status     int
	Code       ErrorCode
	MessageKey string
}

func (e *APIError) Error() string {
	return e.MessageKey
}

func (e *APIError) Message(lang string) string {
	return Translate(lang, e.MessageKey)
}

func NewAPIError(status int, code ErrorCode, messageKey string) *APIError {
	return &APIError{Status: status, Code: code, MessageKey: messageKey}
}

func InvalidParamError(messageKey string) *APIError {
	return NewAPIError(http.StatusBadRequest, CodeCommonInvalidParam, messageKey)
}

func InternalError(messageKey string) *APIError {
	return NewAPIError(http.StatusInternalServerError, CodeInternalError, messageKey)
}

func SetTranslator(t Translator) {
	if t == nil {
		return
	}
	translatorMu.Lock()
	defer translatorMu.Unlock()
	translator = t
}

func Translate(lang, key string) string {
	translatorMu.RLock()
	defer translatorMu.RUnlock()
	return translator.Translate(lang, key)
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

func Write[T any](w http.ResponseWriter, r *http.Request, status int, success bool, code ErrorCode, message string, data T) {
	body := APIResponse[T]{
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

func Success[T any](w http.ResponseWriter, r *http.Request, status int, data T, messageKey string) {
	lang := ResolveLang(r)
	Write(w, r, status, true, CodeOK, Translate(lang, messageKey), data)
}

func Error(w http.ResponseWriter, r *http.Request, status int, code ErrorCode, messageKey string) {
	lang := ResolveLang(r)
	Write[any](w, r, status, false, code, Translate(lang, messageKey), nil)
}

func WriteErrorFrom(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		Error(w, r, apiErr.Status, apiErr.Code, apiErr.MessageKey)
		return
	}

	Error(w, r, http.StatusInternalServerError, CodeInternalError, "error.internal")
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
