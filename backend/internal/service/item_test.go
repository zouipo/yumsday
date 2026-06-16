package service

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	invalidItemID      = int64(-1)
	invalidItemGroupID = int64(-1)
	invalidICID        = int64(-1)

	group1 = model.Group{ID: 1, Name: "Family"}
	group2 = model.Group{ID: 2, Name: "Friends"}

	itemCategory1             = model.ItemCategory{ID: 1, Name: "PANTRY", GroupID: group1.ID}
	itemCategory2             = model.ItemCategory{ID: 2, Name: "BEVERAGE", GroupID: group2.ID}
	itemCategory3             = model.ItemCategory{ID: 3, Name: "DESSERT", GroupID: group1.ID}
	itemCategoryUncategorized = model.ItemCategory{ID: 4, Name: "UNCATEGORIZED", GroupID: group1.ID}

	recipe1 = model.Recipe{ID: 1, Name: "Grilled Chicken"}
	recipe2 = model.Recipe{ID: 2, Name: "Tomato Soup"}

	items = []model.Item{
		{
			ID:                 1,
			Name:               "Flour",
			Description:        new("All-purpose flour"),
			AverageMarketPrice: new(2.50),
			UnitType:           enum.Weight,
			GroupID:            group1.ID,
			ItemCategory:       itemCategory1,
		},
		{
			ID:                 2,
			Name:               "Rice",
			Description:        new("White rice"),
			AverageMarketPrice: new(1.80),
			UnitType:           enum.Weight,
			GroupID:            group1.ID,
			ItemCategory:       itemCategory1,
		},
		{
			ID:                 3,
			Name:               "Water",
			Description:        nil,
			AverageMarketPrice: nil,
			UnitType:           enum.Volume,
			GroupID:            group2.ID,
			ItemCategory:       itemCategory2,
		},
		{
			ID:                 4,
			Name:               "Pepper",
			Description:        nil,
			AverageMarketPrice: new(1.20),
			UnitType:           enum.Weight,
			GroupID:            group1.ID,
			ItemCategory:       itemCategory1,
		},
		{
			ID:                 5,
			Name:               "Olive Oil",
			Description:        new("Extra virgin olive oil"),
			AverageMarketPrice: nil,
			UnitType:           enum.Volume,
			GroupID:            group1.ID,
			ItemCategory:       itemCategory1,
		},
		{
			ID:                 6,
			Name:               "Soft drinks",
			Description:        new("Beverage with no alcohol, usually carbonated"),
			AverageMarketPrice: new(4.00),
			UnitType:           enum.Volume,
			GroupID:            group2.ID,
			ItemCategory:       itemCategory2,
		},
	}
)

type MockItemRepository struct {
	items           []model.Item
	nextID          int64
	getBygroupIDErr error
	getByIDErr      error
	getByNameErr    error
	createErr       error
	updateErr       error
	deleteErr       error
}

func NewMockItemRepository() *MockItemRepository {
	return &MockItemRepository{
		items:  make([]model.Item, 0),
		nextID: 1,
	}
}

/*** MOCK ITEM REPOSITORY (itemRepositoryInterface implementation) ***/

func (m *MockItemRepository) GetByGroupID(groupID int64, sortKey string, desc bool) ([]model.Item, error) {
	if m.getBygroupIDErr != nil {
		return nil, m.getBygroupIDErr
	}

	result := make([]model.Item, 0)
	sortKeys := []string{"name", "averagemarketprice", "unittype.value", "itemcategory.name"}
	if !slices.Contains(sortKeys, strings.ToLower(sortKey)) {
		return nil, customErrors.NewInvalidParamsError([]string{sortKey}, nil)
	}

	for _, item := range m.items {
		if item.GroupID == groupID {
			result = append(result, item)
		}
	}

	return utils.SortSliceByFieldName(result, sortKey, desc), nil
}

func (m *MockItemRepository) GetByID(id int64) (*model.Item, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.items {
		if m.items[i].ID == id {
			return &m.items[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("items", "id", nil)
}

func (m *MockItemRepository) GetByName(name string, desc bool) ([]model.Item, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}

	result := make([]model.Item, 0)
	for _, item := range m.items {
		if item.Name == name {
			result = append(result, item)
		}
	}

	return utils.SortSliceByFieldName(result, "Name", desc), nil
}

func (m *MockItemRepository) Create(item *model.Item) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}

	itemCopy := *item
	itemCopy.ID = m.nextID
	m.nextID++
	m.items = append(m.items, itemCopy)

	return itemCopy.ID, nil
}

func (m *MockItemRepository) Update(item *model.Item) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	for i := range m.items {
		if m.items[i].ID == item.ID {
			// GroupID is not updated
			groupID := m.items[i].GroupID
			m.items[i] = *item
			m.items[i].GroupID = groupID
			return nil
		}
	}

	return customErrors.NewNotFoundError("items", "id", nil)
}

func (m *MockItemRepository) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	for i := range m.items {
		if m.items[i].ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}

	return customErrors.NewNotFoundError("items", "id", nil)
}

/*** MOCK SERVICES ***/

type MockRecipeServiceForItem struct {
	recipes      []model.Recipe
	getByItemErr error
}

// Use by Delete to check if item is used by any recipe
func (m *MockRecipeServiceForItem) GetByItemID(_ int64, _ bool) ([]model.Recipe, error) {
	if m.getByItemErr != nil {
		return nil, m.getByItemErr
	}

	return m.recipes, nil
}

type MockGroceryServiceForItem struct {
	hasItem bool
	err     error
}

// Use by Delete to check if item is used in any grocery
func (m *MockGroceryServiceForItem) HasItem(_ int64) (bool, error) {
	if m.err != nil {
		return false, m.err
	}

	return m.hasItem, nil
}

type MockGroupServiceForItem struct {
	groups     []model.Group
	getByIDErr error
}

// Use by Create and Update to check if group exists
func (m *MockGroupServiceForItem) GetByID(id int64) (*model.Group, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.groups {
		if m.groups[i].ID == id {
			return &m.groups[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("groups", "id", nil)
}

type MockItemCategoryServiceForItem struct {
	itemCategories []model.ItemCategory
	getByIDErr     error
}

// Use by Create and Update to check if item category exists
func (m *MockItemCategoryServiceForItem) GetByID(id int64) (*model.ItemCategory, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.itemCategories {
		if m.itemCategories[i].ID == id {
			return &m.itemCategories[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("item_categories", "id", nil)
}

func (m *MockItemCategoryServiceForItem) GetByNameAndGroupID(name string, groupID int64) (*model.ItemCategory, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	name = strings.ToLower(name)

	for i := range m.itemCategories {
		if strings.ToLower(m.itemCategories[i].Name) == name && m.itemCategories[i].GroupID == groupID {
			return &m.itemCategories[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("item_categories", "name, group_id", nil)
}

/*** HELPER ***/
// Set up test data
func setUpDataTestItem() *MockItemRepository {
	mockRepo := NewMockItemRepository()
	mockRepo.items = append(mockRepo.items, items...)
	mockRepo.nextID = int64(len(mockRepo.items) + 1)
	return mockRepo
}

// Initialize services
func setUpItemCategoryServiceData() *MockItemCategoryServiceForItem {
	return &MockItemCategoryServiceForItem{
		itemCategories: []model.ItemCategory{
			itemCategoryUncategorized,
			itemCategory1,
			itemCategory2,
			itemCategory3,
		},
	}
}

func setUpGroupServiceData() *MockGroupServiceForItem {
	return &MockGroupServiceForItem{
		groups: []model.Group{
			group1,
			group2,
		},
	}
}

func newItemServiceForTest(
	repo *MockItemRepository,
	recipeService *MockRecipeServiceForItem,
	groceryService *MockGroceryServiceForItem,
	groupService *MockGroupServiceForItem,
	itemCategoryService *MockItemCategoryServiceForItem,
) *ItemService {
	return NewItemService(repo, recipeService, groceryService, groupService, itemCategoryService)
}

// getByGroupID returns a slice of items by groupID, sorted by ID
func getByGroupID(id int64, sortKey string, desc bool) []model.Item {
	result := make([]model.Item, 0)
	for _, item := range items {
		if item.GroupID == id {
			result = append(result, item)
		}
	}

	return utils.SortSliceByFieldName(result, sortKey, desc)
}

func compareSlicesItems(s1, s2 []model.Item) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if !reflect.DeepEqual(s1[i], s2[i]) {
			return false
		}
	}

	return true
}

/*** CONSTRUCTOR TEST ***/

func TestNewItemService(t *testing.T) {
	mockRepo := NewMockItemRepository()

	service := NewItemService(
		mockRepo,
		&MockRecipeServiceForItem{},
		&MockGroceryServiceForItem{},
		&MockGroupServiceForItem{},
		&MockItemCategoryServiceForItem{},
	)

	if service == nil {
		t.Fatal("NewItemService() returned nil")
	}

	if service.repo == nil {
		t.Fatal("NewItemService() repo is nil")
	}

	if service.recipeService == nil {
		t.Fatal("NewItemService() recipeService is nil")
	}

	if service.groceryService == nil {
		t.Fatal("NewItemService() groceryService is nil")
	}

	if service.groupService == nil {
		t.Fatal("NewItemService() groupService is nil")
	}

	if service.itemCategoryService == nil {
		t.Fatal("NewItemService() itemCategoryService is nil")
	}
}

/*** READ OPERATIONS TESTS ***/

func TestGetByGroupID(t *testing.T) {
	m := setUpDataTestItem()
	groupService := setUpGroupServiceData()
	s := newItemServiceForTest(
		m,
		&MockRecipeServiceForItem{},
		&MockGroceryServiceForItem{},
		groupService,
		&MockItemCategoryServiceForItem{},
	)

	tests := []struct {
		name        string
		groupID     int64
		sort        string
		descending  bool
		expected    []model.Item
		groupErr    error
		repoErr     error
		expectedErr error
	}{
		{
			name:       "Sort by name asc",
			groupID:    group1.ID,
			sort:       "Name",
			descending: false,
			expected:   getByGroupID(group1.ID, "Name", false),
		},
		{
			name:       "Sort by category desc",
			groupID:    group2.ID,
			sort:       "ItemCategory.Name",
			descending: true,
			expected:   getByGroupID(group2.ID, "ItemCategory.Name", true),
		},
		{
			name:       "Sort by average market price asc",
			groupID:    group1.ID,
			sort:       "AverageMarketPrice",
			descending: false,
			expected:   getByGroupID(group1.ID, "AverageMarketPrice", false),
		},
		{
			name:       "Sort by unit type asc",
			groupID:    group2.ID,
			sort:       "UnitType.value",
			descending: false,
			expected:   getByGroupID(group2.ID, "UnitType.value", false),
		},
		{
			name:        "Group not found",
			groupID:     invalidItemGroupID,
			sort:        "Name",
			expectedErr: customErrors.NewNotFoundError("groups", "id", nil),
		},
		{
			name:        "Group service error",
			groupID:     group1.ID,
			sort:        "Name",
			groupErr:    customErrors.NewInternalError("failed to fetch groups", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch groups", nil),
		},
		{
			name:        "Invalid sort parameter",
			groupID:     group1.ID,
			sort:        "unknown_field",
			expectedErr: customErrors.NewInvalidParamsError([]string{"unknown_field"}, nil),
		},
		{
			name:        "Repository error",
			groupID:     group1.ID,
			sort:        "Name",
			repoErr:     customErrors.NewInternalError("failed to fetch items", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch items", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupService.getByIDErr = tt.groupErr
			m.getBygroupIDErr = tt.repoErr

			actual, err := s.GetByGroupID(tt.groupID, tt.sort, tt.descending)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Errorf("GetByGroupID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Errorf("GetByGroupID() expected nil items on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByGroupID() unexpected error = %v", err)
			}

			if len(actual) != len(tt.expected) {
				t.Errorf("GetByGroupID() returned %d items, expected %d", len(actual), len(tt.expected))
			}

			if !compareSlicesItems(actual, tt.expected) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestGetItemByID(t *testing.T) {
	m := setUpDataTestItem()
	s := newItemServiceForTest(
		m,
		&MockRecipeServiceForItem{},
		&MockGroceryServiceForItem{},
		&MockGroupServiceForItem{},
		&MockItemCategoryServiceForItem{},
	)

	tests := []struct {
		name        string
		itemID      int64
		repoErr     error
		expected    *model.Item
		expectedErr error
	}{
		{
			name:     "Existing ID",
			itemID:   items[0].ID,
			expected: &items[0],
		},
		{
			name:        "Non existing ID",
			itemID:      invalidItemID,
			expectedErr: customErrors.NewNotFoundError("items", "id", nil),
		},
		{
			name:        "Repository error",
			itemID:      items[0].ID,
			repoErr:     customErrors.NewInternalError("failed to fetch items", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch items", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getByIDErr = tt.repoErr

			actual, err := s.GetByID(tt.itemID)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Errorf("GetByID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Errorf("GetByID() expected nil item on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("GetByID() item mismatch: got %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestGetByName(t *testing.T) {
	m := setUpDataTestItem()
	s := newItemServiceForTest(
		m,
		&MockRecipeServiceForItem{},
		&MockGroceryServiceForItem{},
		&MockGroupServiceForItem{},
		&MockItemCategoryServiceForItem{},
	)

	tests := []struct {
		name        string
		itemName    string
		repoErr     error
		expected    []model.Item
		expectedErr error
	}{
		{
			name:     "Existing name",
			itemName: items[0].Name,
			expected: []model.Item{items[0]},
		},
		{
			name:     "Unknown name returns empty slice",
			itemName: "Non existing",
			expected: []model.Item{},
		},
		{
			name:        "Empty name validation",
			itemName:    "",
			expectedErr: customErrors.NewNotFoundError("items", "name", nil),
		},
		{
			name:        "Repository error",
			itemName:    items[0].Name,
			repoErr:     customErrors.NewInternalError("failed to fetch items", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch items", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getByNameErr = tt.repoErr

			actual, err := s.GetByName(tt.itemName, false)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Errorf("GetByName() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Errorf("GetByName() expected nil items on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByName() unexpected error = %v", err)
			}

			if !compareSlicesItems(actual, tt.expected) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestGetRecipes(t *testing.T) {
	tests := []struct {
		name            string
		itemID          int64
		descending      bool
		expectedRecipes []model.Recipe
		recipeErr       error
		expectedErr     error
	}{
		{
			name:            "Existing recipes",
			itemID:          items[0].ID,
			expectedRecipes: []model.Recipe{recipe1, recipe2},
		},
		{
			name:            "Existing recipes descending",
			itemID:          items[0].ID,
			descending:      true,
			expectedRecipes: []model.Recipe{recipe2, recipe1},
		},
		{
			name:        "Repository error",
			itemID:      items[0].ID,
			recipeErr:   customErrors.NewInternalError("failed to fetch recipes", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch recipes", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipeService := &MockRecipeServiceForItem{
				recipes:      tt.expectedRecipes,
				getByItemErr: tt.recipeErr,
			}

			s := newItemServiceForTest(
				setUpDataTestItem(),
				recipeService,
				&MockGroceryServiceForItem{},
				&MockGroupServiceForItem{},
				&MockItemCategoryServiceForItem{},
			)

			actual, err := s.GetRecipesByID(tt.itemID, tt.descending)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Errorf("GetRecipesByID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Errorf("GetRecipesByID() expected nil recipes on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetRecipesByID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expectedRecipes) {
				t.Errorf("GetRecipesByID() recipes mismatch: got %v, want %v", actual, tt.expectedRecipes)
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name         string
		item         *model.Item
		icErr        error
		groupErr     error
		repoErr      error
		expectedErr  error
		expectedItem *model.Item
	}{
		{
			name: "valid item",
			item: &model.Item{
				Name:               "Pasta",
				Description:        new("Dry spaghetti"),
				AverageMarketPrice: new(2.10),
				UnitType:           enum.Weight,
				GroupID:            itemCategory1.GroupID,
				ItemCategory:       itemCategory1,
			},
			expectedItem: &model.Item{
				ID:                 int64(len(items) + 1),
				Name:               "Pasta",
				Description:        new("Dry spaghetti"),
				AverageMarketPrice: new(2.10),
				UnitType:           enum.Weight,
				GroupID:            itemCategory1.GroupID,
				ItemCategory:       itemCategory1,
			},
		},
		{
			name: "valid item, no category provided",
			item: &model.Item{
				Name:               "Sparkling Water",
				Description:        new("Carbonated water"),
				AverageMarketPrice: new(1.30),
				UnitType:           enum.Volume,
				GroupID:            group1.ID,
			},
			expectedItem: &model.Item{
				ID:                 int64(len(items) + 1),
				Name:               "Sparkling Water",
				Description:        new("Carbonated water"),
				AverageMarketPrice: new(1.30),
				UnitType:           enum.Volume,
				GroupID:            group1.ID,
				ItemCategory:       itemCategoryUncategorized,
			},
		},
		{
			name: "Validation error empty name",
			item: &model.Item{
				Name:         "",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"name"}, nil),
		},
		{
			name: "Validation error empty unit type",
			item: &model.Item{
				Name:         "Pasta",
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"unit type"}, nil),
		},
		{
			name: "Validation error unit type and name",
			item: &model.Item{
				Name:         "",
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"name", "unit type"}, nil),
		},
		{
			name: "Item category not found",
			item: &model.Item{
				Name:     "Pasta",
				UnitType: enum.Weight,
				GroupID:  group1.ID,
				ItemCategory: model.ItemCategory{
					ID: int64(invalidIcID),
				},
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "Item category service error",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			icErr:       customErrors.NewInternalError("failed to fetch item categories", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch item categories", nil),
		},
		{
			name: "Item category belongs to another group",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory2,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "Group not found",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      int64(invalidItemGroupID),
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewConflictError("Group", "group must exists", nil),
		},
		{
			name: "Group service error",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			groupErr:    customErrors.NewInternalError("failed to fetch groups", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch groups", nil),
		},
		{
			name: "Repository error",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			repoErr:     customErrors.NewInternalError("failed to create item", nil),
			expectedErr: customErrors.NewInternalError("failed to create item", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setUpDataTestItem()
			repo.createErr = tt.repoErr

			groupService := setUpGroupServiceData()
			groupService.getByIDErr = tt.groupErr

			itemCategoryService := setUpItemCategoryServiceData()
			itemCategoryService.getByIDErr = tt.icErr

			s := newItemServiceForTest(
				repo,
				&MockRecipeServiceForItem{},
				&MockGroceryServiceForItem{},
				groupService,
				itemCategoryService,
			)

			id, err := s.Create(tt.item)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Errorf("Create() error = %v, want %v", err, tt.expectedErr)
				}
				if id != 0 {
					t.Errorf("Create() expected ID 0 on error, got %d", id)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error = %v", err)
			}

			if id == 0 {
				t.Errorf("Create() returned ID 0, expected non-zero")
			}

			if len(repo.items) != int(tt.expectedItem.ID) {
				t.Errorf("Create() items count = %d, want %d", len(repo.items), int(tt.expectedItem.ID))
			}

			newItem, err := repo.GetByID(id)
			if err != nil {
				t.Fatalf("GetByID() after Update() error = %v", err)
			}

			if !reflect.DeepEqual(newItem, tt.expectedItem) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expectedItem, newItem)
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name         string
		item         *model.Item
		icErr        error
		groupErr     error
		repoErr      error
		expectedErr  error
		expectedItem *model.Item
	}{
		{
			name: "valid update with valid new category",
			item: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",                  // modified
				Description:        new("Premium all-purpose flour"), // modified
				AverageMarketPrice: new(3.20),                        // modified
				UnitType:           enum.Volume,                      // modified
				GroupID:            items[0].GroupID,
				ItemCategory:       itemCategory3, // modified, but still coherent with the group
			},
			expectedItem: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Volume,
				GroupID:            items[0].GroupID,
				ItemCategory:       itemCategory3,
			},
		},
		{
			name: "Same update, ignore group absence",
			item: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				ItemCategory:       itemCategory3,
			},
			expectedItem: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				GroupID:            items[0].GroupID,
				ItemCategory:       itemCategory3,
			},
		},
		{
			name: "Same update, ignore invalid group",
			item: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				GroupID:            invalidICGroupID,
				ItemCategory:       itemCategory3,
			},
			expectedItem: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				GroupID:            items[0].GroupID,
				ItemCategory:       itemCategory3,
			},
		},
		{
			name: "No category provided, uses Uncategorized",
			item: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
			},
			expectedItem: &model.Item{
				ID:                 items[0].ID,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				GroupID:            items[0].GroupID,
				ItemCategory:       itemCategoryUncategorized,
			},
		},
		{
			name: "Validation name error",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"name"}, nil),
		},
		{
			name: "Validation unit type error",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "Flour Premium",
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"unit type"}, nil),
		},
		{
			name: "Category not found",
			item: &model.Item{
				ID:       items[0].ID,
				Name:     "Flour Premium",
				UnitType: enum.Weight,
				GroupID:  group1.ID,
				ItemCategory: model.ItemCategory{
					ID: invalidICID,
				},
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "Category belongs to another group",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "Flour Premium",
				UnitType:     enum.Weight,
				ItemCategory: itemCategory2,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "Repository error",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "Flour Premium",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory1,
			},
			repoErr:     customErrors.NewInternalError("failed to update item", nil),
			expectedErr: customErrors.NewInternalError("failed to update item", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setUpDataTestItem()
			repo.updateErr = tt.repoErr

			groupService := setUpGroupServiceData()
			groupService.getByIDErr = tt.groupErr

			itemCategoryService := setUpItemCategoryServiceData()
			itemCategoryService.getByIDErr = tt.icErr

			s := newItemServiceForTest(
				repo,
				&MockRecipeServiceForItem{},
				&MockGroceryServiceForItem{},
				groupService,
				itemCategoryService,
			)

			err := s.Update(tt.item)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Errorf("Update() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() unexpected error = %v", err)
			}

			updated, err := repo.GetByID(tt.item.ID)
			if err != nil {
				t.Fatalf("GetByID() after Update() error = %v", err)
			}

			if !reflect.DeepEqual(updated, tt.expectedItem) {
				t.Errorf("Items should be equal: expected %v, got %v", tt.expectedItem, updated)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name        string
		itemID      int64
		recipes     []model.Recipe
		recipeErr   error
		hasItem     bool
		groceryErr  error
		repoErr     error
		expectedErr error
	}{
		{
			name:   "Success",
			itemID: items[0].ID,
		},
		{
			name:        "Used in recipe",
			itemID:      items[0].ID,
			recipes:     []model.Recipe{recipe1},
			expectedErr: customErrors.NewConflictError("Item", "can't delete item used by recipes", nil),
		},
		{
			name:        "Recipe service error",
			itemID:      items[0].ID,
			recipeErr:   customErrors.NewInternalError("failed to fetch recipes", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch recipes", nil),
		},
		{
			name:        "Used in grocery",
			itemID:      items[0].ID,
			hasItem:     true,
			expectedErr: customErrors.NewConflictError("Item", "can't delete item used in groceries", nil),
		},
		{
			name:        "Grocery service error",
			itemID:      items[0].ID,
			groceryErr:  customErrors.NewInternalError("failed to check grocery item", nil),
			expectedErr: customErrors.NewInternalError("failed to check grocery item", nil),
		},
		{
			name:        "Repository error",
			itemID:      items[0].ID,
			repoErr:     customErrors.NewInternalError("failed to delete item", nil),
			expectedErr: customErrors.NewInternalError("failed to delete item", nil),
		},
		{
			name:        "Not found",
			itemID:      invalidItemID,
			expectedErr: customErrors.NewNotFoundError("items", "id", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setUpDataTestItem()
			repo.deleteErr = tt.repoErr

			recipeService := &MockRecipeServiceForItem{
				recipes:      tt.recipes,
				getByItemErr: tt.recipeErr,
			}

			groceryService := &MockGroceryServiceForItem{
				hasItem: tt.hasItem,
				err:     tt.groceryErr,
			}

			s := newItemServiceForTest(
				repo,
				recipeService,
				groceryService,
				&MockGroupServiceForItem{},
				&MockItemCategoryServiceForItem{},
			)

			err := s.Delete(tt.itemID)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("Delete() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Delete() unexpected error = %v", err)
			}

			deleted, err := repo.GetByID(tt.itemID)
			if err == nil || deleted != nil {
				t.Fatalf("Delete() item with ID %d should be deleted", tt.itemID)
			}
		})
	}
}

func TestCheckSimpleFields(t *testing.T) {
	tests := []struct {
		name        string
		item        *model.Item
		expectedErr error
	}{
		{
			name: "valid",
			item: &model.Item{Name: "Flour", UnitType: enum.Weight},
		},
		{
			name:        "empty name",
			item:        &model.Item{Name: "", UnitType: enum.Weight},
			expectedErr: customErrors.NewInvalidParamsError([]string{"name"}, nil),
		},
		{
			name:        "empty unit type",
			item:        &model.Item{Name: "Flour"},
			expectedErr: customErrors.NewInvalidParamsError([]string{"unit type"}, nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSimpleFields(tt.item)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("checkSimpleFields() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("checkSimpleFields() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateItem(t *testing.T) {
	tests := []struct {
		name        string
		item        *model.Item
		icErr       error
		groupErr    error
		expectedErr error
	}{
		{
			name: "valid new item",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
		},
		{
			name: "valid existing item",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      items[0].GroupID,
				ItemCategory: items[0].ItemCategory,
			},
		},
		{
			name: "invalid name",
			item: &model.Item{
				Name:         "",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"name"}, nil),
		},
		{
			name: "invalid unit type",
			item: &model.Item{
				Name:         "Pasta",
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewInvalidParamsError([]string{"unit type"}, nil),
		},
		{
			name: "item category not found",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: model.ItemCategory{ID: invalidICID},
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "item category internal error",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			icErr:       customErrors.NewInternalError("failed to fetch item categories", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch item categories", nil),
		},
		{
			name: "item category belongs to another group, new item",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory2,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "item category belongs to another group, existing item",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      items[0].GroupID,
				ItemCategory: itemCategory2,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "group not found, new item",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      int64(invalidItemGroupID),
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewConflictError("Group", "group must exists", nil),
		},
		{
			name: "group not found, existing item",
			item: &model.Item{
				ID:           items[0].ID,
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      int64(invalidItemGroupID),
				ItemCategory: items[0].ItemCategory,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "group internal error, new item",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			groupErr:    customErrors.NewInternalError("failed to fetch groups", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch groups", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setUpDataTestItem()
			groupService := setUpGroupServiceData()
			groupService.getByIDErr = tt.groupErr

			itemCategoryService := setUpItemCategoryServiceData()
			itemCategoryService.getByIDErr = tt.icErr

			s := newItemServiceForTest(
				repo,
				&MockRecipeServiceForItem{},
				&MockGroceryServiceForItem{},
				groupService,
				itemCategoryService,
			)

			err := s.validateItem(tt.item)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("validateItem() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("validateItem() unexpected error = %v", err)
			}
		})
	}
}
