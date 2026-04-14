package service

import (
	"reflect"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	icID             = 1
	invalidIcID      = -1
	invalidICGroupID = int64(-1)
	invalidICName    = "UNKNOWN"
)

type MockItemCategoryRepository struct {
	itemCategories         []model.ItemCategory
	getByIDErr             error
	getByNameAndGroupIDErr error
}

func NewMockItemCategoryRepository() *MockItemCategoryRepository {
	return &MockItemCategoryRepository{
		itemCategories: make([]model.ItemCategory, 0),
	}
}

func (m *MockItemCategoryRepository) GetByID(id int64) (*model.ItemCategory, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.itemCategories {
		if m.itemCategories[i].ID == id {
			return &m.itemCategories[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("ItemCategory", "items.id", nil)
}

func (m *MockItemCategoryRepository) GetByNameAndGroupID(name string, groupID int64) (*model.ItemCategory, error) {
	if m.getByNameAndGroupIDErr != nil {
		return nil, m.getByNameAndGroupIDErr
	}

	for i := range m.itemCategories {
		if m.itemCategories[i].Name == name && m.itemCategories[i].GroupID == groupID {
			return &m.itemCategories[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("ItemCategory", "items.name, items.group_id", nil)
}

func setUpDataTestIC() *MockItemCategoryRepository {
	mockRepo := NewMockItemCategoryRepository()
	mockRepo.itemCategories = append(mockRepo.itemCategories, model.ItemCategory{
		ID:      int64(icID),
		Name:    "FRUITS",
		GroupID: 1,
	})
	mockRepo.itemCategories = append(mockRepo.itemCategories, model.ItemCategory{
		ID:      int64(icID + 1),
		Name:    "VEGETABLES",
		GroupID: 2,
	})
	mockRepo.itemCategories = append(mockRepo.itemCategories, model.ItemCategory{
		ID:      int64(icID + 2),
		Name:    "DAIRY",
		GroupID: 1,
	})

	return mockRepo
}

func TestNewItemCategoryService(t *testing.T) {
	mockRepo := &MockItemCategoryRepository{}

	service := NewItemCategoryService(mockRepo)

	if service == nil {
		t.Fatal("NewItemCategoryService returned nil")
	}

	if service.repo == nil {
		t.Fatal("NewItemCategoryService repo is nil")
	}
}

func TestGetItemCategoryByID(t *testing.T) {
	m := setUpDataTestIC()
	s := NewItemCategoryService(m)

	tests := []struct {
		name        string
		icID        int64
		expected    *model.ItemCategory
		err         error
		expectedErr error
	}{
		{
			name:     "Existing ID",
			icID:     int64(icID),
			expected: &utils.SortSliceByFieldName(m.itemCategories, "ID", false)[0],
		},
		{
			name:        "Non existing ID",
			icID:        int64(invalidIcID),
			expected:    nil,
			expectedErr: customErrors.NewNotFoundError("ItemCategory", "items.id", nil),
		},
		{
			name:        "Repository error",
			icID:        int64(icID),
			expected:    nil,
			expectedErr: customErrors.NewInternalError("failed to fetch item categories", nil),
			err:         customErrors.NewInternalError("failed to fetch item categories", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getByIDErr = tt.err

			actual, err := s.GetByID(tt.icID)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("GetByID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Fatalf("GetByID() expected nil item category on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("GetByID() item categories mismatch: got %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestGetItemCategoryByNameAndGroupID(t *testing.T) {
	m := setUpDataTestIC()
	s := NewItemCategoryService(m)

	repoError := customErrors.NewInternalError("failed to fetch item categories", nil)

	firstItem := utils.SortSliceByFieldName(m.itemCategories, "ID", false)[0]
	secondItem := utils.SortSliceByFieldName(m.itemCategories, "ID", false)[1]

	tests := []struct {
		name        string
		icName      string
		groupID     int64
		expected    *model.ItemCategory
		err         error
		expectedErr error
	}{
		{
			name:     "Existing name and group ID",
			icName:   firstItem.Name,
			groupID:  firstItem.GroupID,
			expected: &firstItem,
		},
		{
			name:        "Non existing name and group ID",
			icName:      invalidICName,
			groupID:     invalidICGroupID,
			expected:    nil,
			expectedErr: customErrors.NewNotFoundError("ItemCategory", "items.name, items.group_id", nil),
		},
		{
			name:        "Existing name, bad group ID",
			icName:      firstItem.Name,
			groupID:     secondItem.GroupID,
			expected:    nil,
			expectedErr: customErrors.NewNotFoundError("ItemCategory", "items.name, items.group_id", nil),
		},
		{
			name:        "Repository error",
			icName:      m.itemCategories[0].Name,
			groupID:     m.itemCategories[0].GroupID,
			expected:    nil,
			expectedErr: repoError,
			err:         repoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getByNameAndGroupIDErr = tt.err

			actual, err := s.GetByNameAndGroupID(tt.icName, tt.groupID)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("GetByNameAndGroupID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Fatalf("GetByNameAndGroupID() expected nil item category on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByNameAndGroupID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("GetByNameAndGroupID() item categories mismatch: got %v, want %v", actual, tt.expected)
			}
		})
	}
}
