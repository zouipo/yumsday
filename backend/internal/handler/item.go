package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zouipo/yumsday/backend/internal/constant"
	"github.com/zouipo/yumsday/backend/internal/dto"
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
	mux.HandleFunc("POST "+prefix, h.createItem)
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

// createItem add a new item to the database
// @Summary Create a new item
// @Description Create a new item with the provided details
// @Tags item
// @Accept json
// @ Produce json
// @Param item body dto.ItemDto true "New Item data"
// @Success 201 {object} map[string]int "Returns the new item ID"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/item [post]
func (h *ItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	var newItemDto dto.ItemDto
	err := json.NewDecoder(r.Body).Decode(&newItemDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := mapper.ToItem(&newItemDto)
	id, err := h.itemService.Create(item)
	if err != nil {
		if appErr, ok := errors.AsType[customErrors.AppError](err); ok {
			http.Error(w, err.Error(), appErr.HTTPStatus())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constant.CONTENT_TYPE_HEADER, constant.CONTENT_TYPE_VALUE)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, id)
}

// updateItem updates an existing item
// @Summary Update item details
// @Description Update the details of an existing item
// @Tags item
// @Accept json
// @Produce json
// @Param item body dto.ItemDto true "Item data to update"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Item not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/item [put]
func (h *ItemHandler) updateItem(w http.ResponseWriter, r *http.Request) {
	var itemDto dto.ItemDto
	if err := json.NewDecoder(r.Body).Decode(&itemDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := mapper.ToItem(&itemDto)
	if err := h.itemService.Update(item); err != nil {
		if appErr, ok := errors.AsType[customErrors.AppError](err); ok {
			http.Error(w, err.Error(), appErr.HTTPStatus())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constant.CONTENT_TYPE_HEADER, constant.CONTENT_TYPE_VALUE)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteItem delete an item
// @Summary Delete an item
// @Description Delete the item with the specified ID
// @Tags item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Item not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/item/{id} [delete]
func (h *ItemHandler) deleteItem(w http.ResponseWriter, r *http.Request) {
	err := h.itemService.Delete(r.Context().Value("id").(int64))

	if err != nil {
		if appErr, ok := errors.AsType[customErrors.AppError](err); ok {
			http.Error(w, err.Error(), appErr.HTTPStatus())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constant.CONTENT_TYPE_HEADER, constant.CONTENT_TYPE_VALUE)
	w.WriteHeader(http.StatusNoContent)
}
