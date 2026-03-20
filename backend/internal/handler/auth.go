package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/dto"
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
// @Param credentials body dto.LoginDto true "Login credentials"
// @Success 302 {string} string "Redirect to home page"
// @Failure 400 {string} string "Missing username or password"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (h *AuthHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginDto
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		http.Error(w, "missing username or password", http.StatusBadRequest)
		return
	}

	session := r.Context().Value(ctx.SessionCtxKey{}).(*model.Session)
	err = h.s.Authenticate(session, loginReq.Username, loginReq.Password)
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
