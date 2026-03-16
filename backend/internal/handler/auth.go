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

// GetLogin godoc
// @Summary Get login page
// @Description Retrieve the login page
// @Tags auth
// @Produce html
// @Success 200 {string} string "Login page HTML"
// @Failure 500 {string} string "Internal server error"
// @Router /login [get]
func (h *AuthHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login page !"))
}

// PostLogin godoc
// @Summary Authenticate user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 302 {string} string "Redirect to home page"
// @Failure 400 {string} string "Missing username or password"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
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

// PostLogout godoc
// @Summary Logout user
// @Description Logout the authenticated user
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirect to login page"
// @Failure 500 {string} string "Internal server error"
// @Router /logout [post]
func (h *AuthHandler) postLogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(ctx.SessionCtxKey{}).(*model.Session)
	err := h.s.Logout(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
