package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/zouipo/yumsday/backend/internal/constant"
	"github.com/zouipo/yumsday/backend/internal/dto"
	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/mapper"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

var (
	ItemCategory1 = &model.ItemCategory{
		ID:      1,
		Name:    "GRAINS AND PASTA",
		GroupID: 1,
	}
	ItemCategory2 = &model.ItemCategory{
		ID:      2,
		Name:    "VEGETABLES",
		GroupID: 2,
	}
	ItemCategory3 = &model.ItemCategory{
		ID:      3,
		Name:    "SPICES AND CONDIMENTS",
		GroupID: 1,
	}
	testItem1 = &model.Item{
		ID:                 1,
		Name:               "Flour",
		Description:        new("All-purpose flour"),
		AverageMarketPrice: new(2.50),
		UnitType:           enum.Weight,
		ItemCategory:       *ItemCategory1,
		GroupID:            1,
	}
	testItem2 = &model.Item{
		ID:                 2,
		Name:               "Onions",
		Description:        new("Yellow onions"),
		AverageMarketPrice: new(1.50),
		UnitType:           enum.Weight,
		ItemCategory:       *ItemCategory2,
		GroupID:            2,
	}
	testItem3 = &model.Item{
		ID:           3,
		Name:         "Olive Oil",
		Description:  new("Extra virgin olive oil"),
		UnitType:     enum.Volume,
		ItemCategory: *ItemCategory3,
		GroupID:      1,
	}

	validItemID = 1

	invalidItemID   = -1
	invalidItemName = "psd"
)

/*** MOCK SERVICE ***/
type MockItemService struct {
	items             []model.Item
	nextID            int64
	GetByIDErr        error
	GetByNameErr      error
	GetRecipesByIDErr error
	CreateErr         error
	UpdateErr         error
	DeleteErr         error
}

func NewMockItemService() *MockItemService {
	return &MockItemService{
		items:  make([]model.Item, 0),
		nextID: 1,
	}
}

func (s *MockItemService) GetByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error) {
	return nil, nil
}

func (m *MockItemService) GetByID(id int64) (*model.Item, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}

	for i := range m.items {
		if m.items[i].ID == id {
			return &m.items[i], nil
		}
	}
	return nil, customErrors.NewNotFoundError("items", strconv.FormatInt(id, 10), errors.New(userNotFoundErr))
}

func (s *MockItemService) GetByName(name string, descending bool) ([]model.Item, error) {
	return make([]model.Item, 0), nil
}

func (s *MockItemService) GetRecipesByID(id int64, descending bool) ([]model.Recipe, error) {
	return make([]model.Recipe, 0), nil
}

func (s *MockItemService) Create(item *model.Item) (int64, error) {
	return 0, nil
}

func (s *MockItemService) Update(item *model.Item) error {
	return nil
}

func (s *MockItemService) Delete(id int64) error {
	return nil
}

/*** HELPER FUNCTIONS ***/

func (m *MockItemService) addItem(item *model.Item) {
	item.ID = m.nextID
	m.nextID++
	m.items = append(m.items, *item)
}

func setupItemTestData() *MockItemService {
	mockService := NewMockItemService()

	mockService.addItem(testItem1)
	mockService.addItem(testItem2)
	mockService.addItem(testItem3)

	return mockService
}

/*** TEST CONSTRUCTOR ***/

func TestNewItemHandler(t *testing.T) {
	mockService := NewMockItemService()
	handler := NewItemHandler(mockService)

	if handler == nil {
		t.Fatal("expected non-nil handler")

		if handler.itemService != mockService {
			t.Error("handler itemService does not match the provided service")
		}
	}
}

/*** READ OPERATIONS TESTS ***/

func TestGetByID(t *testing.T) {
	mockService := setupItemTestData()
	handler := NewItemHandler(mockService)

	tests := []struct {
		name     string
		itemID   int64
		expected *dto.ItemDto
		code     int64
		err      error
	}{
		{
			name:     "Success",
			itemID:   1,
			code:     http.StatusOK,
			expected: mapper.ToItemDto(testItem1),
		},
		{
			name:   "Non existing ID",
			itemID: int64(invalidItemID),
			code:   http.StatusNotFound,
			err:    customErrors.NewNotFoundError("items", "id", nil),
		},
		{
			name:   "Internal server error",
			itemID: int64(invalidItemID),
			code:   http.StatusInternalServerError,
			err:    customErrors.NewInternalError("failed to fetch items", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				mockService.GetByIDErr = tt.err
			}

			r := httptest.NewRequest(http.MethodGet, "/item/"+strconv.FormatInt(tt.itemID, 10), nil)
			// Add the ID to the context as the middleware would do
			ctx := context.WithValue(r.Context(), "id", tt.itemID)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.getItemById(w, r)

			if w.Code != int(tt.code) {
				t.Errorf("expected status %d instead of %d", tt.code, w.Code)
			}

			// If success
			if tt.err == nil {
				contentType := w.Header().Get(constant.CONTENT_TYPE_HEADER)
				if contentType != constant.CONTENT_TYPE_VALUE {
					t.Errorf("expected content type %s instead of %s", constant.CONTENT_TYPE_VALUE, contentType)
				}

				var actual dto.ItemDto
				err := json.NewDecoder(w.Body).Decode(&actual)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if reflect.DeepEqual(actual, tt.expected) {
					t.Errorf("Actual item %v mismatched expected item %v", actual, tt.expected)
				}
			}
		})
	}
}
