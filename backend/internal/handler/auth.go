package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/constant"
	"github.com/zouipo/yumsday/backend/internal/ctx"
	"github.com/zouipo/yumsday/backend/internal/dto"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/mapper"
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
	mux.HandleFunc("POST /login", h.postLogin)
	mux.HandleFunc("POST /logout", h.postLogout)
}

// PostLogin godoc
// @Summary Authenticate user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginDto true "Login credentials"
// @Success 200 {string} string "Login successful"
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
	user, err := h.s.Authenticate(session, loginReq.Username, loginReq.Password)
	if err != nil {
		if appErr, ok := errors.AsType[customErrors.AppError](err); ok {
			http.Error(w, err.Error(), appErr.HTTPStatus())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constant.CONTENT_TYPE_HEADER, constant.CONTENT_TYPE_VALUE)
	if err = json.NewEncoder(w).Encode(mapper.ToUserDtoNoPassword(user)); err != nil {
		http.Error(w, "Failed to serialize user", http.StatusInternalServerError)
		return
	}
}

// PostLogout godoc
// @Summary Logout user
// @Description Logout the authenticated user
// @Tags auth
// @Produce json
// @Success 204 {string} string "Logout successful"
// @Failure 500 {string} string "Internal server error"
// @Router /logout [post]
func (h *AuthHandler) postLogout(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(ctx.SessionCtxKey{}).(*model.Session)
	err := h.s.Logout(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
