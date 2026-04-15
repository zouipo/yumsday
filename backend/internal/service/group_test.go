package service

import (
	"reflect"
	"testing"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	groupID        = 1
	invalidGroupID = -1
)

type MockGroupRepository struct {
	groups     []model.Group
	getByIDErr error
}

func NewMockGroupRepository() *MockGroupRepository {
	return &MockGroupRepository{
		groups: make([]model.Group, 0),
	}
}

func (m *MockGroupRepository) GetByID(id int64) (*model.Group, error) {
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

func setUpDataTestGroup() *MockGroupRepository {
	mockRepo := NewMockGroupRepository()
	mockRepo.groups = append(mockRepo.groups, model.Group{
		ID:        int64(groupID),
		Name:      "Family",
		ImageURL:  new("/static/images/family.jpg"),
		CreatedAt: time.Now(),
		Members: []model.GroupMember{
			{
				UserID:   1,
				Admin:    true,
				JoinedAt: time.Now().AddDate(0, 0, -1),
			},
			{
				UserID:   2,
				Admin:    false,
				JoinedAt: time.Now().AddDate(0, 0, -1),
			},
		},
	})
	mockRepo.groups = append(mockRepo.groups, model.Group{
		ID:        int64(groupID + 1),
		Name:      "Friends",
		ImageURL:  new("/static/images/friends.jpg"),
		CreatedAt: time.Now(),
		Members:   []model.GroupMember{},
	})

	return mockRepo
}

func TestNewGroupService(t *testing.T) {
	mockRepo := &MockGroupRepository{}

	service := NewGroupService(mockRepo)

	if service == nil {
		t.Fatal("NewItemCategoryService returned nil")
	}

	if service.repo == nil {
		t.Fatal("NewItemCategoryService repo is nil")
	}
}

func TestGetGroupByID(t *testing.T) {
	m := setUpDataTestGroup()
	s := NewGroupService(m)

	tests := []struct {
		name        string
		groupID     int64
		expected    *model.Group
		err         error
		expectedErr error
	}{
		{
			name:     "Existing ID, multiple group members",
			groupID:  int64(groupID),
			expected: &utils.SortSliceByFieldName(m.groups, "ID", false)[0],
		},
		{
			name:     "Existing ID, no group members",
			groupID:  int64(groupID + 1),
			expected: &utils.SortSliceByFieldName(m.groups, "ID", false)[1],
		},
		{
			name:        "Non existing ID",
			groupID:     int64(invalidGroupID),
			expected:    nil,
			expectedErr: customErrors.NewNotFoundError("Group", "groups.id", nil),
		},
		{
			name:        "Repository error",
			groupID:     int64(groupID),
			expected:    nil,
			expectedErr: customErrors.NewInternalError("failed to fetch groups", nil),
			err:         customErrors.NewInternalError("failed to fetch groups", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.getByIDErr = tt.err

			actual, err := s.GetByID(tt.groupID)

			if tt.expectedErr != nil {
				if !utils.CompareErrors(err, tt.expectedErr) {
					t.Fatalf("GetByID() error = %v, want %v", err, tt.expectedErr)
				}
				if actual != nil {
					t.Fatalf("GetByID() expected nil group on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByID() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("GetByID() groups mismatch: got %v, want %v", actual, tt.expected)
			}
		})
	}
}
