package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"vibe-golang-template/internal/model"
	"vibe-golang-template/internal/service"
	"vibe-golang-template/pkg/response"
)

type API struct {
	users *service.UserService
}

func NewAPI(users *service.UserService) *API {
	return &API{users: users}
}

func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", a.healthz)
	mux.HandleFunc("/api/v1/users", a.usersRoute)
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) usersRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.listUsers(w)
	case http.MethodPost:
		a.createUser(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) listUsers(w http.ResponseWriter) {
	response.JSON(w, http.StatusOK, map[string]any{"items": a.users.ListUsers()})
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	created, err := a.users.CreateUser(model.User{Name: in.Name, Email: in.Email})
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserInput) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusCreated, created)
}
