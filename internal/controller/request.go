package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"vibe-golang-template/pkg/response"
)

const (
	defaultPage  = 1
	defaultSize  = 20
	maxPageSize  = 100
	defaultLimit = 20
	maxLimit     = 100
)

type PageParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type CursorParams struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

func decodeJSONBody(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return response.InvalidParamError("error.invalid_json_body")
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return response.InvalidParamError("error.json_single_object")
	}

	return nil
}

func parsePageParams(r *http.Request) (PageParams, error) {
	page := defaultPage
	size := defaultSize

	if raw := strings.TrimSpace(r.URL.Query().Get("page")); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 {
			return PageParams{}, response.InvalidParamError("error.page_invalid")
		}
		page = v
	}

	if raw := strings.TrimSpace(r.URL.Query().Get("size")); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 || v > maxPageSize {
			return PageParams{}, response.InvalidParamError("error.size_invalid")
		}
		size = v
	}

	return PageParams{Page: page, Size: size}, nil
}

func parseCursorParams(r *http.Request) (CursorParams, error) {
	limit := defaultLimit
	cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))

	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 || v > maxLimit {
			return CursorParams{}, response.InvalidParamError("error.limit_invalid")
		}
		limit = v
	}

	if cursor != "" && strings.TrimSpace(r.URL.Query().Get("page")) != "" {
		return CursorParams{}, response.InvalidParamError("error.cursor_page_conflict")
	}

	return CursorParams{Cursor: cursor, Limit: limit}, nil
}
