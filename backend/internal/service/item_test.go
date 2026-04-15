package service

import (
	"reflect"
	"sort"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	itemID             = int64(1)
	invalidItemID      = int64(-1)
	invalidItemGroupID = int64(-1)
	invalidICID        = int64(-1)

	group1 = model.Group{
		ID:   1,
		Name: "Family",
	}
	group2 = model.Group{
		ID:   2,
		Name: "Friends",
	}

	itemCategory1 = model.ItemCategory{
		ID:      1,
		Name:    "PANTRY",
		GroupID: group1.ID,
	}

	itemCategory2 = model.ItemCategory{
		ID:      2,
		Name:    "BEVERAGE",
		GroupID: group2.ID,
	}

	itemCategoryUncategorized1 = model.ItemCategory{
		ID:      3,
		Name:    "Uncategorized",
		GroupID: group1.ID,
	}

	itemCategoryUncategorized2 = model.ItemCategory{
		ID:      4,
		Name:    "Uncategorized",
		GroupID: group2.ID,
	}

	// ID 0 allows validateItem(GetByID(0)) to pass before Create/Update fallback logic runs.
	itemCategoryZeroGroup1 = model.ItemCategory{
		ID:      0,
		Name:    "TEMP ZERO",
		GroupID: group1.ID,
	}

	itemCategoryZeroGroup2 = model.ItemCategory{
		ID:      0,
		Name:    "TEMP ZERO",
		GroupID: group2.ID,
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
	lastSortKey     string
	lastDescending  bool
}

func NewMockItemRepository() *MockItemRepository {
	return &MockItemRepository{
		items:  make([]model.Item, 0),
		nextID: 1,
	}
}

/*** MOCK ITEM REPOSITORY (itemRepositoryInterface implementation) ***/

func (m *MockItemRepository) GetByGroupID(groupID int64, sortKey string, descending bool) ([]model.Item, error) {
	if m.getBygroupIDErr != nil {
		return nil, m.getBygroupIDErr
	}

	m.lastSortKey = sortKey
	m.lastDescending = descending

	result := make([]model.Item, 0)
	for _, item := range m.items {
		if item.GroupID == groupID {
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		left := result[i]
		right := result[j]

		switch sortKey {
		case "i.average_market_price":
			// Treat nil prices as -1 so items without a price can still be ordered deterministically.
			leftPrice := -1.0
			rightPrice := -1.0
			if left.AverageMarketPrice != nil {
				leftPrice = *left.AverageMarketPrice
			}
			if right.AverageMarketPrice != nil {
				rightPrice = *right.AverageMarketPrice
			}
			if descending {
				return leftPrice > rightPrice
			}
			return leftPrice < rightPrice
		case "i.unit_type":
			if descending {
				return left.UnitType.String() > right.UnitType.String()
			}
			return left.UnitType.String() < right.UnitType.String()
		case "ic.name":
			if descending {
				return left.ItemCategory.Name > right.ItemCategory.Name
			}
			return left.ItemCategory.Name < right.ItemCategory.Name
		default:
			if descending {
				return left.Name > right.Name
			}
			return left.Name < right.Name
		}
	})

	return result, nil
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

	return nil, customErrors.NewNotFoundError("Item", "items.id", nil)
}

func (m *MockItemRepository) GetByName(name string) ([]model.Item, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}

	result := make([]model.Item, 0)
	for _, item := range m.items {
		if item.Name == name {
			result = append(result, item)
		}
	}

	return result, nil
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
			m.items[i] = *item
			return nil
		}
	}

	return customErrors.NewNotFoundError("Item", "items.id", nil)
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

	return customErrors.NewNotFoundError("Item", "items.id", nil)
}

/*** MOCK SERVICES ***/

type MockRecipeServiceForItem struct {
	recipes      []model.Recipe
	getByItemErr error
}

// Use by Delete to check if item is used by any recipe
func (m *MockRecipeServiceForItem) GetByItemID(_ int64) ([]model.Recipe, error) {
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

	return nil, customErrors.NewNotFoundError("Group", "groups.id", nil)
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

	return nil, customErrors.NewNotFoundError("ItemCategory", "item_categories.id", nil)
}

func (m *MockItemCategoryServiceForItem) GetByNameAndGroupID(name string, groupID int64) (*model.ItemCategory, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.itemCategories {
		if m.itemCategories[i].Name == name && m.itemCategories[i].GroupID == groupID {
			return &m.itemCategories[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("ItemCategory", "item_categories.name AND item_categories.group_id", nil)
}

/*** HELPER ***/

func setUpDataTestItem() *MockItemRepository {
	mockRepo := NewMockItemRepository()

	mockRepo.items = append(mockRepo.items, model.Item{
		ID:                 1,
		Name:               "Flour",
		Description:        new("All-purpose flour"),
		AverageMarketPrice: new(2.50),
		UnitType:           enum.Weight,
		GroupID:            group1.ID,
		ItemCategory:       itemCategory1,
	})
	mockRepo.items = append(mockRepo.items, model.Item{
		ID:                 2,
		Name:               "Rice",
		Description:        new("White rice"),
		AverageMarketPrice: new(1.80),
		UnitType:           enum.Weight,
		GroupID:            group1.ID,
		ItemCategory:       itemCategory1,
	})
	mockRepo.items = append(mockRepo.items, model.Item{
		ID:                 3,
		Name:               "Water",
		Description:        nil,
		AverageMarketPrice: nil,
		UnitType:           enum.Volume,
		GroupID:            group2.ID,
		ItemCategory:       itemCategory2,
	})
	mockRepo.items = append(mockRepo.items, model.Item{
		ID:                 4,
		Name:               "Pepper",
		Description:        nil,
		AverageMarketPrice: new(1.20),
		UnitType:           enum.Weight,
		GroupID:            group1.ID,
		ItemCategory:       itemCategory1,
	})
	mockRepo.items = append(mockRepo.items, model.Item{
		ID:                 5,
		Name:               "Olive Oil",
		Description:        new("Extra virgin olive oil"),
		AverageMarketPrice: nil,
		UnitType:           enum.Volume,
		GroupID:            group1.ID,
		ItemCategory:       itemCategory1,
	})
	mockRepo.items = append(mockRepo.items, model.Item{
		ID:                 6,
		Name:               "Soft drinks",
		Description:        new("Beverage with no alcohol, usually carbonated"),
		AverageMarketPrice: new(4.00),
		UnitType:           enum.Volume,
		GroupID:            group2.ID,
		ItemCategory:       itemCategory2,
	})

	mockRepo.nextID = int64(len(mockRepo.items) + 1)

	return mockRepo
}

func setUpItemCategoryServiceData() *MockItemCategoryServiceForItem {
	return &MockItemCategoryServiceForItem{
		itemCategories: []model.ItemCategory{
			itemCategoryUncategorized1,
			itemCategoryUncategorized2,
			itemCategory1,
			itemCategory2,
			itemCategoryZeroGroup1,
			itemCategoryZeroGroup2,
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

func TestGetByGroupID(t *testing.T) {
	m := setUpDataTestItem()
	s := newItemServiceForTest(
		m,
		&MockRecipeServiceForItem{},
		&MockGroceryServiceForItem{},
		&MockGroupServiceForItem{},
		&MockItemCategoryServiceForItem{},
	)

	tests := []struct {
		name            string
		groupID         int64
		sort            string
		descending      bool
		expectedSortKey string
		expectedLen     int
		repoErr         error
		expectedErr     error
		expectedFirstID int64
	}{
		{
			name:            "Sort by name asc",
			groupID:         group1.ID,
			sort:            "name",
			descending:      false,
			expectedSortKey: "i.name",
			expectedLen:     4,
			expectedFirstID: 1,
		},
		{
			name:            "Sort by category desc",
			groupID:         group2.ID,
			sort:            "category",
			descending:      true,
			expectedSortKey: "ic.name",
			expectedLen:     2,
			expectedFirstID: 0,
		},
		{
			name:            "Sort by average market price",
			groupID:         group1.ID,
			sort:            "average_market_price",
			descending:      false,
			expectedSortKey: "i.average_market_price",
			expectedLen:     4,
			expectedFirstID: 5,
		},
		{
			name:            "Sort by unit type",
			groupID:         group2.ID,
			sort:            "unit_type",
			descending:      false,
			expectedSortKey: "i.unit_type",
			expectedLen:     2,
			expectedFirstID: 3,
		},
		{
			name:        "Invalid sort parameter",
			groupID:     group1.ID,
			sort:        "unknown_field",
			expectedErr: customErrors.NewInvalidParamsError("unknown_field", nil),
		},
		{
			name:        "Repository error",
			groupID:     group1.ID,
			sort:        "name",
			repoErr:     customErrors.NewInternalError("failed to fetch items", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch items", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getBygroupIDErr = tt.repoErr

			actual, err := s.GetByGroupID(tt.groupID, tt.sort, tt.descending)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("GetByGroupID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Fatalf("GetByGroupID() expected nil items on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByGroupID() unexpected error = %v", err)
			}

			if m.lastSortKey != tt.expectedSortKey {
				t.Fatalf("GetByGroupID() sort key = %s, want %s", m.lastSortKey, tt.expectedSortKey)
			}

			if m.lastDescending != tt.descending {
				t.Fatalf("GetByGroupID() descending = %v, want %v", m.lastDescending, tt.descending)
			}

			if len(actual) != tt.expectedLen {
				t.Fatalf("GetByGroupID() returned %d items, expected %d", len(actual), tt.expectedLen)
			}

			if tt.expectedFirstID > 0 && len(actual) > 0 && actual[0].ID != tt.expectedFirstID {
				t.Fatalf("GetByGroupID() first item ID = %d, want %d", actual[0].ID, tt.expectedFirstID)
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
			name:   "Existing ID",
			itemID: itemID,
			expected: &model.Item{
				ID:                 1,
				Name:               "Flour",
				Description:        new("All-purpose flour"),
				AverageMarketPrice: new(2.50),
				UnitType:           enum.Weight,
				GroupID:            group1.ID,
				ItemCategory: model.ItemCategory{
					ID:      itemCategory1.ID,
					Name:    "PANTRY",
					GroupID: group1.ID,
				},
			},
		},
		{
			name:        "Non existing ID",
			itemID:      invalidItemID,
			expectedErr: customErrors.NewNotFoundError("Item", "items.id", nil),
		},
		{
			name:        "Repository error",
			itemID:      itemID,
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
					t.Fatalf("GetByID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Fatalf("GetByID() expected nil item on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("GetByID() item mismatch: got %v, want %v", actual, tt.expected)
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
		expectedLen int
		expectedErr error
	}{
		{
			name:        "Existing name",
			itemName:    "Flour",
			expectedLen: 1,
		},
		{
			name:        "Unknown name returns empty slice",
			itemName:    "Non existing",
			expectedLen: 0,
		},
		{
			name:        "Empty name validation",
			itemName:    "",
			expectedErr: customErrors.NewNotFoundError("Item", "", nil),
		},
		{
			name:        "Repository error",
			itemName:    "Flour",
			repoErr:     customErrors.NewInternalError("failed to fetch items", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch items", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getByNameErr = tt.repoErr

			actual, err := s.GetByName(tt.itemName)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("GetByName() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Fatalf("GetByName() expected nil items on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByName() unexpected error = %v", err)
			}

			if len(actual) != tt.expectedLen {
				t.Fatalf("GetByName() returned %d items, expected %d", len(actual), tt.expectedLen)
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name                 string
		item                 *model.Item
		icErr                error
		groupErr             error
		repoErr              error
		expectedErr          error
		expectedCreatedItems int
		expectedCategoryID   int64
	}{
		{
			name: "Success",
			item: &model.Item{
				Name:               "Pasta",
				Description:        new("Dry spaghetti"),
				AverageMarketPrice: new(2.10),
				UnitType:           enum.Weight,
				GroupID:            itemCategory1.GroupID,
				ItemCategory:       itemCategory1,
			},
			expectedCreatedItems: 7,
			expectedCategoryID:   itemCategory1.ID,
		},
		{
			name: "No category provided uses Uncategorized",
			item: &model.Item{
				Name:               "Sparkling Water",
				Description:        new("Carbonated water"),
				AverageMarketPrice: new(1.30),
				UnitType:           enum.Volume,
				GroupID:            group1.ID,
				ItemCategory:       model.ItemCategory{},
			},
			expectedCreatedItems: 7,
			expectedCategoryID:   itemCategoryUncategorized1.ID,
		},
		{
			name: "Validation error empty name",
			item: &model.Item{
				Name:         "",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewValidationError("name", "item must have a name", nil),
		},
		{
			name: "Validation error empty unit type",
			item: &model.Item{
				Name:         "Pasta",
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewValidationError("unit type", "item must have a unit type", nil),
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
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
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
					t.Fatalf("Create() error = %v, want %v", err, tt.expectedErr)
				}
				if id != 0 {
					t.Fatalf("Create() expected ID 0 on error, got %d", id)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create() unexpected error = %v", err)
			}

			if id == 0 {
				t.Fatal("Create() returned ID 0, expected non-zero")
			}

			if len(repo.items) != tt.expectedCreatedItems {
				t.Fatalf("Create() items count = %d, want %d", len(repo.items), tt.expectedCreatedItems)
			}

			if tt.expectedCategoryID > 0 {
				created, getErr := repo.GetByID(id)
				if getErr != nil {
					t.Fatalf("GetByID() after Create() error = %v", getErr)
				}
				if created.ItemCategory.ID != tt.expectedCategoryID {
					t.Fatalf("Create() item category ID = %d, want %d", created.ItemCategory.ID, tt.expectedCategoryID)
				}
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name               string
		item               *model.Item
		icErr              error
		groupErr           error
		repoErr            error
		expectedErr        error
		expectedCategoryID int64
	}{
		{
			name: "Success",
			item: &model.Item{
				ID:                 1,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				GroupID:            itemCategory1.GroupID,
				ItemCategory:       itemCategory1,
			},
			expectedCategoryID: itemCategory1.ID,
		},
		{
			name: "No category provided uses Uncategorized",
			item: &model.Item{
				ID:                 1,
				Name:               "Flour Premium",
				Description:        new("Premium all-purpose flour"),
				AverageMarketPrice: new(3.20),
				UnitType:           enum.Weight,
				GroupID:            group1.ID,
				ItemCategory:       model.ItemCategory{},
			},
			expectedCategoryID: itemCategoryUncategorized1.ID,
		},
		{
			name: "Validation error",
			item: &model.Item{
				ID:           1,
				Name:         "",
				UnitType:     enum.Weight,
				GroupID:      itemCategory1.GroupID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewValidationError("name", "item must have a name", nil),
		},
		{
			name: "Category not found",
			item: &model.Item{
				ID:       1,
				Name:     "Flour Premium",
				UnitType: enum.Weight,
				GroupID:  1,
				ItemCategory: model.ItemCategory{
					ID: invalidICID,
				},
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "Category belongs to another group",
			item: &model.Item{
				ID:           1,
				Name:         "Flour Premium",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory2,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "Group not found",
			item: &model.Item{
				ID:       1,
				Name:     "Flour Premium",
				UnitType: enum.Weight,
				GroupID:  int64(invalidItemGroupID),
				ItemCategory: model.ItemCategory{
					ID: itemCategory1.ID,
				},
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "Repository error",
			item: &model.Item{
				ID:       1,
				Name:     "Flour Premium",
				UnitType: enum.Weight,
				GroupID:  1,
				ItemCategory: model.ItemCategory{
					ID: 1,
				},
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
					t.Fatalf("Update() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Update() unexpected error = %v", err)
			}

			updated, getErr := repo.GetByID(tt.item.ID)
			if getErr != nil {
				t.Fatalf("GetByID() after Update() error = %v", getErr)
			}

			if updated.Name != tt.item.Name {
				t.Fatalf("Update() updated name = %s, want %s", updated.Name, tt.item.Name)
			}

			if tt.expectedCategoryID > 0 && updated.ItemCategory.ID != tt.expectedCategoryID {
				t.Fatalf("Update() item category ID = %d, want %d", updated.ItemCategory.ID, tt.expectedCategoryID)
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
			itemID: 3,
		},
		{
			name:   "Used in recipe",
			itemID: 6,
			recipes: []model.Recipe{
				{ID: 1, Name: "Grilled Chicken"},
			},
			expectedErr: customErrors.NewConflictError("Item", "can't delete item used by recipes", nil),
		},
		{
			name:        "Recipe service error",
			itemID:      6,
			recipeErr:   customErrors.NewInternalError("failed to fetch recipes", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch recipes", nil),
		},
		{
			name:        "Used in grocery",
			itemID:      3,
			hasItem:     true,
			expectedErr: customErrors.NewConflictError("Item", "can't delete item used in groceries", nil),
		},
		{
			name:        "Grocery service error",
			itemID:      3,
			groceryErr:  customErrors.NewInternalError("failed to check grocery item", nil),
			expectedErr: customErrors.NewInternalError("failed to check grocery item", nil),
		},
		{
			name:        "Repository error",
			itemID:      3,
			repoErr:     customErrors.NewInternalError("failed to delete item", nil),
			expectedErr: customErrors.NewInternalError("failed to delete item", nil),
		},
		{
			name:        "Not found",
			itemID:      invalidItemID,
			expectedErr: customErrors.NewNotFoundError("Item", "items.id", nil),
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

			deleted, getErr := repo.GetByID(tt.itemID)
			if getErr == nil || deleted != nil {
				t.Fatalf("Delete() item with ID %d should be deleted", tt.itemID)
			}
		})
	}
}

func TestMapSortKey(t *testing.T) {
	repo := setUpDataTestItem()
	s := newItemServiceForTest(
		repo,
		&MockRecipeServiceForItem{},
		&MockGroceryServiceForItem{},
		&MockGroupServiceForItem{},
		&MockItemCategoryServiceForItem{},
	)

	tests := []struct {
		name        string
		param       string
		expectedKey string
		expectedErr error
	}{
		{name: "default empty", param: "", expectedKey: "i.name"},
		{name: "name", param: "name", expectedKey: "i.name"},
		{name: "average market price", param: "average_market_price", expectedKey: "i.average_market_price"},
		{name: "unit type", param: "unit_type", expectedKey: "i.unit_type"},
		{name: "category", param: "category", expectedKey: "ic.name"},
		{name: "invalid", param: "wrong", expectedErr: customErrors.NewInvalidParamsError("wrong", nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := s.mapSortKey(tt.param)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("mapSortKey() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != "" {
					t.Fatalf("mapSortKey() expected empty key on error, got %q", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("mapSortKey() unexpected error = %v", err)
			}

			if actual != tt.expectedKey {
				t.Fatalf("mapSortKey() key = %q, want %q", actual, tt.expectedKey)
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
			item: &model.Item{
				Name:     "Flour",
				UnitType: enum.Weight,
			},
		},
		{
			name: "empty name",
			item: &model.Item{
				Name:     "",
				UnitType: enum.Weight,
			},
			expectedErr: customErrors.NewValidationError("name", "item must have a name", nil),
		},
		{
			name: "empty unit type",
			item: &model.Item{
				Name: "Flour",
			},
			expectedErr: customErrors.NewValidationError("unit type", "item must have a unit type", nil),
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
			name: "valid",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory1,
			},
		},
		{
			name: "invalid name",
			item: &model.Item{
				Name:         "",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewValidationError("name", "item must have a name", nil),
		},
		{
			name: "invalid unit type",
			item: &model.Item{
				Name:         "Pasta",
				GroupID:      group1.ID,
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewValidationError("unit type", "item must have a unit type", nil),
		},
		{
			name: "item category not found",
			item: &model.Item{
				Name:     "Pasta",
				UnitType: enum.Weight,
				GroupID:  group1.ID,
				ItemCategory: model.ItemCategory{
					ID: invalidICID,
				},
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "item category internal error",
			item: &model.Item{
				Name:     "Pasta",
				UnitType: enum.Weight,
				GroupID:  1,
				ItemCategory: model.ItemCategory{
					ID: 1,
				},
			},
			icErr:       customErrors.NewInternalError("failed to fetch item categories", nil),
			expectedErr: customErrors.NewInternalError("failed to fetch item categories", nil),
		},
		{
			name: "item category belongs to another group",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
				ItemCategory: itemCategory2,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil),
		},
		{
			name: "group not found",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      int64(invalidItemGroupID),
				ItemCategory: itemCategory1,
			},
			expectedErr: customErrors.NewConflictError("ItemCategory", "item category must exists", nil),
		},
		{
			name: "group internal error",
			item: &model.Item{
				Name:         "Pasta",
				UnitType:     enum.Weight,
				GroupID:      group1.ID,
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
