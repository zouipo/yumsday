package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
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

	invalidItemID   = -1
	invalidItemName = "psd"
)

/*** MOCK SERVICE ***/
type MockItemService struct {
	items             []model.Item
	nextID            int64
	getByIDErr        error
	getByNameErr      error
	getRecipesByIDErr error
	createErr         error
	updateErr         error
	deleteErr         error
}

func NewMockItemService() *MockItemService {
	return &MockItemService{
		items:  make([]model.Item, 0),
		nextID: 1,
	}
}

func (m *MockItemService) GetByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error) {
	return nil, nil
}

func (m *MockItemService) GetByID(id int64) (*model.Item, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.items {
		if m.items[i].ID == id {
			return &m.items[i], nil
		}
	}
	return nil, customErrors.NewNotFoundError("items", strconv.FormatInt(id, 10), errors.New(userNotFoundErr))
}

func (m *MockItemService) GetByName(groupID int64, name string, descending bool) ([]model.Item, error) {
	return make([]model.Item, 0), nil
}

func (m *MockItemService) GetRecipesByID(id int64, descending bool) ([]model.Recipe, error) {
	return make([]model.Recipe, 0), nil
}

func (m *MockItemService) Create(item *model.Item) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}
	item.ID = m.nextID
	m.nextID++

	m.items = append(m.items, *item)

	return item.ID, nil
}

func (m *MockItemService) Update(item *model.Item) error {
	return nil
}

func (m *MockItemService) Delete(id int64) error {
	return nil
}

/*** HELPER FUNCTIONS ***/

func (m *MockItemService) addItem(item *model.Item) {
	item.ID = m.nextID
	m.nextID++
	m.items = append(m.items, *item)
}

// setupItemTestData creates a fresh mock service with predefined test items for test independence.
// It is run at the start of each test to ensure a consistent state and avoid test interference.
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
				mockService.getByIDErr = tt.err
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

				if tt.expected == nil || !reflect.DeepEqual(actual, *tt.expected) {
					t.Errorf("Actual item %v mismatched expected item %v", actual, tt.expected)
				}
			}
		})
	}
}

/*** CREATE OPERATIONS TESTS ***/

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name         string
		itemDto      dto.ItemDto
		expectedItem model.Item
		code         int64
		err          error
	}{
		{
			name: "Success - no ID provided",
			itemDto: dto.ItemDto{
				Name:               "Tagliatelle",
				Description:        new("Long & flat pasta"),
				AverageMarketPrice: new(2.50),
				UnitType:           enum.Weight,
				ItemCategory:       *mapper.ToItemCategoryDto(ItemCategory1),
				GroupID:            1,
			},
			code: http.StatusCreated,
			expectedItem: model.Item{
				ID:                 4,
				Name:               "Tagliatelle",
				Description:        new("Long & flat pasta"),
				AverageMarketPrice: new(2.50),
				UnitType:           enum.Weight,
				ItemCategory:       *ItemCategory1,
				GroupID:            1,
			},
		},
		{
			name: "Success - no ID provided & pre-existing fields",
			itemDto: dto.ItemDto{
				Name:               "Flour",
				Description:        new("All-purpose flour"),
				AverageMarketPrice: new(2.50),
				UnitType:           enum.Weight,
				ItemCategory:       *mapper.ToItemCategoryDto(ItemCategory1),
				GroupID:            1,
			},
			code:         http.StatusCreated,
			expectedItem: *testItem1,
		},
		{
			name:         "Success - ignore dto ID",
			itemDto:      *mapper.ToItemDto(testItem1),
			code:         http.StatusCreated,
			expectedItem: *testItem1,
		},
		{
			name:    "Internal server error",
			itemDto: *mapper.ToItemDto(testItem1),
			code:    http.StatusInternalServerError,
			err:     customErrors.NewInternalError("failed to create item", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := setupItemTestData()
			handler := NewItemHandler(mockService)

			if tt.err != nil {
				mockService.createErr = tt.err
			}

			itemsNb := len(mockService.items)

			body, _ := json.Marshal(tt.itemDto)
			r := httptest.NewRequest(http.MethodPost, "/item", bytes.NewReader(body))
			r.Header.Set(constant.CONTENT_TYPE_HEADER, constant.CONTENT_TYPE_VALUE)
			w := httptest.NewRecorder()

			handler.createItem(w, r)

			if w.Code != int(tt.code) {
				t.Errorf("expected status %d instead of %d", tt.code, w.Code)
			}

			// If success
			if tt.err == nil {
				contentType := w.Header().Get(constant.CONTENT_TYPE_HEADER)
				if contentType != constant.CONTENT_TYPE_VALUE {
					t.Errorf("expected content type %s instead of %s", constant.CONTENT_TYPE_VALUE, contentType)
				}

				var result map[string]int
				err := json.NewDecoder(w.Body).Decode(&result)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				var expectedID = int(mockService.nextID - 1)
				if result["id"] != expectedID {
					t.Errorf("expected id %d instead of %d", expectedID, result["id"])
				}

				if itemsNb+1 != len(mockService.items) {
					t.Errorf("expected %d users instead of %d", itemsNb+1, len(mockService.items))
				}

				item, err := mockService.GetByID((int64(result["id"])))
				if err != nil {
					t.Fatalf("failed to retrieve created item: %v", err)
				}

				tt.expectedItem.ID = int64(expectedID)
				if !reflect.DeepEqual(*item, tt.expectedItem) {
					t.Errorf("Actual item %v mismatched the expected %v", *item, tt.expectedItem)
				}
			}
		})
	}
}

func TestCreateItemHandler_DecodeErrors(t *testing.T) {
	mockService := NewMockItemService()
	handler := NewItemHandler(mockService)

	tests := []struct {
		name string
		body string
		code int
		err  string
	}{
		{
			name: "empty body",
			body: "",
			code: http.StatusBadRequest,
			err:  "EOF",
		},
		{
			name: "malformed JSON syntax",
			body: `{"id": 1, "name": "Test Item",`,
			code: http.StatusBadRequest,
			err:  "unexpected EOF",
		},
		{
			name: "wrong type for id",
			body: `{"id": "not-a-number", "name": "Test Item", "group_id": 1}`,
			code: http.StatusBadRequest,
			err:  "cannot unmarshal string into Go struct field ItemDto.id",
		},
		{
			name: "invalid unit_type value",
			body: `{"id": 1, "name": "Test Item", "group_id": 1, "unit_type": "TEST"}`,
			code: http.StatusBadRequest,
			err:  "invalid unit type value: TEST",
		},
		{
			name: "non-string unit_type",
			body: `{"id": 1, "name": "Test Item", "group_id": 1, "unit_type": 42}`,
			code: http.StatusBadRequest,
			err:  "json: cannot unmarshal number into Go struct field ItemDto.unit_type of type string",
		},
		{
			name: "empty string unit_type",
			body: `{"id": 1, "name": "Test Item", "group_id": 1, "unit_type": ""}`,
			code: http.StatusBadRequest,
			err:  "invalid unit type value:",
		},
		{
			name: "top-level array instead of object",
			body: `[1, 2, 3]`,
			code: http.StatusBadRequest,
			err:  "cannot unmarshal array into Go value of type dto.ItemDto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.createItem(rec, req)

			if rec.Code != tt.code {
				t.Errorf("Actual status code = %d, expected %d", rec.Code, tt.code)
			}

			gotBody := strings.TrimSpace(rec.Body.String())
			if !strings.Contains(gotBody, tt.err) {
				t.Errorf("Actual response body = %q, expected it to contain %q", gotBody, tt.err)
			}
		})
	}
}
