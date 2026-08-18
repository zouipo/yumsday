package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/constant"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/mapper"
	"github.com/zouipo/yumsday/backend/internal/middleware"
	"github.com/zouipo/yumsday/backend/internal/service"
)

// ItemHandler handles HTTP requests related to item operations.
type ItemHandler struct {
	itemService service.ItemServiceInterface
}

// NewItemHandler constructs a new ItemHandler with the provided ItemService.
func NewItemHandler(itemService service.ItemServiceInterface) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.Handle("GET "+prefix+"/{id}", middleware.IntPathValues("id")(http.HandlerFunc(h.getItemById)))
}

// GetItemByID fetchs an item by its ID
// @Summary Get item by ID
// @Description Get an item by its ID
// @Tags item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} dto.ItemDto
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Item not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/item/{id} [get]
func (h *ItemHandler) getItemById(w http.ResponseWriter, r *http.Request) {
	item, err := h.itemService.GetByID(r.Context().Value("id").(int64))
	if err != nil {
		if appErr, ok := errors.AsType[customErrors.AppError](err); ok {
			http.Error(w, err.Error(), appErr.HTTPStatus())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constant.CONTENT_TYPE_HEADER, constant.CONTENT_TYPE_VALUE)
	if err = json.NewEncoder(w).Encode(mapper.ToItemDto(item)); err != nil {
		http.Error(w, customErrors.SERIALIZE_USER_ERROR, http.StatusInternalServerError)
		return
	}
}
