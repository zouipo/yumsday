package handler

import (
	"errors"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/service"
)

type AuthHandler struct {
	s service.AuthServiceInterface
}

func NewAuthHandler(s service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		s: s,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /login", h.getLogin)
	mux.HandleFunc("POST /login", h.postLogin)
	mux.HandleFunc("POST /logout", h.postLogout)
}

// getLogin serves the login page.
func (h *AuthHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login page !"))
}

// postLogin handles the login form submission.
func (h *AuthHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "missing username or password", http.StatusBadRequest)
		return
	}

	session := r.Context().Value(ctx.SessionCtxKey{}).(*model.Session)
	err := h.s.Authenticate(session, username, password)
	if err != nil {
		if appErr, ok := errors.AsType[*customErrors.AppError](err); ok {
			http.Error(w, err.Error(), appErr.StatusCode)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) postLogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(ctx.SessionCtxKey{}).(*model.Session)
	err := h.s.Logout(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
