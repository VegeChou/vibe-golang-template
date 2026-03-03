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
	mux.HandleFunc("/_rules/health", a.healthz)
	mux.HandleFunc("/api/v1/users", a.usersRoute)
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	lang := response.ResolveLang(r)
	if r.Method != http.MethodGet {
		response.Error(
			w,
			r,
			http.StatusMethodNotAllowed,
			response.CodeCommonInvalidParam,
			response.Message(lang, "method not allowed", "请求方法不允许"),
		)
		return
	}
	response.Success(w, r, http.StatusOK, map[string]string{"status": "UP"}, response.Message(lang, "success", "成功"))
}

func (a *API) usersRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.listUsers(w, r)
	case http.MethodPost:
		a.createUser(w, r)
	default:
		lang := response.ResolveLang(r)
		response.Error(
			w,
			r,
			http.StatusMethodNotAllowed,
			response.CodeCommonInvalidParam,
			response.Message(lang, "method not allowed", "请求方法不允许"),
		)
	}
}

func (a *API) listUsers(w http.ResponseWriter, r *http.Request) {
	lang := response.ResolveLang(r)
	response.Success(
		w,
		r,
		http.StatusOK,
		map[string]any{"items": a.users.ListUsers()},
		response.Message(lang, "success", "成功"),
	)
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		lang := response.ResolveLang(r)
		response.Error(
			w,
			r,
			http.StatusBadRequest,
			response.CodeCommonInvalidParam,
			response.Message(lang, "invalid json body", "JSON 请求体无效"),
		)
		return
	}

	created, err := a.users.CreateUser(model.User{Name: in.Name, Email: in.Email})
	if err != nil {
		lang := response.ResolveLang(r)
		if errors.Is(err, service.ErrInvalidUserInput) {
			response.Error(
				w,
				r,
				http.StatusBadRequest,
				response.CodeCommonInvalidParam,
				response.Message(lang, "invalid user input", "用户输入不合法"),
			)
			return
		}
		response.Error(
			w,
			r,
			http.StatusInternalServerError,
			response.CodeInternalError,
			response.Message(lang, "internal error", "系统内部错误"),
		)
		return
	}

	lang := response.ResolveLang(r)
	response.Success(w, r, http.StatusCreated, created, response.Message(lang, "success", "成功"))
}
