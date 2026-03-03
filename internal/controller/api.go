package controller

import (
	"errors"
	"math"
	"net/http"

	"vibe-golang-template/internal/model"
	"vibe-golang-template/internal/service"
	"vibe-golang-template/pkg/response"
)

type API struct {
	users *service.UserService
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type pageUsersData struct {
	Items      []model.User `json:"items"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	Total      int          `json:"total"`
	TotalPages int          `json:"totalPages"`
	HasNext    bool         `json:"hasNext"`
}

type routeFunc func(http.ResponseWriter, *http.Request) error

func NewAPI(users *service.UserService) *API {
	return &API{users: users}
}

func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", a.withError(a.healthz))
	mux.HandleFunc("/_rules/health", a.withError(a.healthz))
	mux.HandleFunc("/api/v1/users", a.withError(a.usersRoute))
}

func (a *API) withError(fn routeFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			response.WriteErrorFrom(w, r, err)
		}
	}
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return response.NewAPIError(http.StatusMethodNotAllowed, response.CodeCommonInvalidParam, "error.method_not_allowed")
	}

	response.Success(w, r, http.StatusOK, map[string]string{"status": "UP"}, "common.success")
	return nil
}

func (a *API) usersRoute(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return a.listUsers(w, r)
	case http.MethodPost:
		return a.createUser(w, r)
	default:
		return response.NewAPIError(http.StatusMethodNotAllowed, response.CodeCommonInvalidParam, "error.method_not_allowed")
	}
}

func (a *API) listUsers(w http.ResponseWriter, r *http.Request) error {
	pageParams, err := parsePageParams(r)
	if err != nil {
		return err
	}
	if _, err := parseCursorParams(r); err != nil {
		return err
	}

	allUsers := a.users.ListUsers()
	total := len(allUsers)
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageParams.Size)))
	}

	start := (pageParams.Page - 1) * pageParams.Size
	end := start + pageParams.Size
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	items := make([]model.User, 0)
	if start < end {
		items = allUsers[start:end]
	}

	response.Success(
		w,
		r,
		http.StatusOK,
		pageUsersData{
			Items:      items,
			Page:       pageParams.Page,
			Size:       pageParams.Size,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    pageParams.Page < totalPages,
		},
		"common.success",
	)
	return nil
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) error {
	var in createUserRequest
	if err := decodeJSONBody(r, &in); err != nil {
		return err
	}

	created, err := a.users.CreateUser(model.User{Name: in.Name, Email: in.Email})
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserInput) {
			return response.InvalidParamError("error.invalid_user_input")
		}
		return response.InternalError("error.internal")
	}

	response.Success(w, r, http.StatusCreated, created, "common.success")
	return nil
}
