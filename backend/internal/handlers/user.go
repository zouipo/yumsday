package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/dtos"
	"github.com/zouipo/yumsday/backend/internal/mappers"
	"github.com/zouipo/yumsday/backend/internal/services"
)

// UserHandler handles HTTP requests related to user operations.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler constructs a new UserHandler with the provided UserService.
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers the user-related routes on the provided ServeMux with the given prefix.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("GET "+prefix, h.getUsers)
	mux.HandleFunc("GET "+prefix+"/{id}", h.getUserByID)
	mux.HandleFunc("POST "+prefix, h.createUser)
	mux.HandleFunc("PUT "+prefix, h.updateUser)
	mux.HandleFunc("PATCH "+prefix+"/{id}/role", h.updateUserAdminRole)
	mux.HandleFunc("PATCH "+prefix+"/{id}/password", h.updateUserPassword)
	mux.HandleFunc("DELETE "+prefix+"/{id}", h.deleteUser)
}

// GetUsers godoc
// @Summary Get users
// @Description Get all users or a user by username
// @Tags user
// @Accept json
// @Produce json
// @Param username query string false "Username to filter by"
// @Success 200 {array} dtos.UserDto
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /user [get]
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		h.getAllUsers(w)
		return
	}

	usernames := queryParams["username"]
	if len(usernames) == 1 {
		h.getByUsername(w, usernames[0])
		return
	}

	http.Error(w, "Missing or invalid query parameters", http.StatusBadRequest)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.UserDto
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [get]
func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetByID(r.Context().Value("id").(int64))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mappers.ToUserDtoNoPassword(user))
	if err != nil {
		http.Error(w, "Failed to serialize user", http.StatusInternalServerError)
		return
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags user
// @Accept json
// @Produce json
// @Param user body dtos.NewUserDto true "New User Data"
// @Success 201 {object} map[string]int "Returns the new user ID"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /user [post]
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var newUserDto dtos.NewUserDto
	err := json.NewDecoder(r.Body).Decode(&newUserDto)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := mappers.FromNewUserDtoToUser(&newUserDto)
	id, err := h.userService.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, id)
}

// UpdateUser godoc
// @Summary Update user details
// @Description Update the details of an existing user
// @Tags user
// @Accept json
// @Produce json
// @Param user body dtos.UserDto true "User Data to Update"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /user [put]
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.UserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := mappers.FromUserDtoToUser(&userDto)
	err = h.userService.Update(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserAdminRole godoc
// @Summary Update user admin role
// @Description Update the admin role status for a specific user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param role body map[string]bool true "Admin Role Status"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id}/admin [patch]
func (h *UserHandler) updateUserAdminRole(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(int64)

	err := json.NewDecoder(r.Body).Decode(&adminRolePayload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateAdminRole(userID, adminRolePayload.AppAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserPassword godoc
// @Summary Update user password
// @Description Update the password for a specific user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param password body map[string]string true "Old and New Passwords"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id}/password [patch]
func (h *UserHandler) updateUserPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(int64)

	err := json.NewDecoder(r.Body).Decode(&passwordPayload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdatePassword(userID, passwordPayload.OldPassword, passwordPayload.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete the user with the specified ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [delete]
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.userService.Delete(r.Context().Value("id").(int64))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

/*** PRIVATE STRUCT ***/
var adminRolePayload struct {
	AppAdmin bool `json:"app_admin"`
}

var passwordPayload struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

/*** NON-HANDLER PRIVATE METHODS ***/
// getAllUsers retrieves all users and writes them to the response.
func (h *UserHandler) getAllUsers(w http.ResponseWriter) {
	users, err := h.userService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mappers.MapList(users, mappers.ToUserDtoNoPassword))
	if err != nil {
		http.Error(w, "Failed to serialize users", http.StatusInternalServerError)
		return
	}
}

// getByUsername retrieves a user by username and writes it to the response as an array.
func (h *UserHandler) getByUsername(w http.ResponseWriter, username string) {
	user, err := h.userService.GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return as an array with one user to match the array response of the original handler getUsers.
	users := []*dtos.UserDto{mappers.ToUserDtoNoPassword(user)}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Failed to serialize user", http.StatusInternalServerError)
		return
	}
}
