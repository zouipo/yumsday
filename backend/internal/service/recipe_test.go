package service

import (
	"context"
	"reflect"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

type MockRecipeRepository struct {
	recipes      []model.Recipe
	getByItemErr error
}

func NewMockRecipeRepository() *MockRecipeRepository {
	return &MockRecipeRepository{
		recipes: make([]model.Recipe, 0),
	}
}

func (m *MockRecipeRepository) GetByID(_ int64) (model.Recipe, error) {
	return model.Recipe{}, nil
}

func (m *MockRecipeRepository) GetByGroupID(_ int64) ([]model.Recipe, error) {
	return nil, nil
}

func (m *MockRecipeRepository) GetByItemID(itemID int64) ([]model.Recipe, error) {
	if m.getByItemErr != nil {
		return nil, m.getByItemErr
	}

	return m.recipes, nil
}

func (m *MockRecipeRepository) Create(_ context.Context, _ *model.Recipe) (int64, error) {
	return 0, nil
}

func (m *MockRecipeRepository) Update(_ *model.Recipe) error {
	return nil
}

func (m *MockRecipeRepository) Delete(_ int64) error {
	return nil
}

/*** MOCK SERVICES ***/
type MockIngredientForRecipe struct {
	validateIngError error
}

func (m *MockIngredientForRecipe) validateIngredient(ing model.Ingredient) error {
	if m.validateIngError != nil {
		return m.validateIngError
	}

	return nil
}

type MockCategoryForRecipe struct {
	validateCatError error
}

func (m *MockCategoryForRecipe) validateRecipeCategory(cat model.RecipeCategory) error {
	if m.validateCatError != nil {
		return m.validateCatError
	}

	return nil
}

type MockGroupForRecipe struct {
	GetByIDError error
	group        *model.Group
}

func (m *MockGroupForRecipe) GetByID(_ int64) (*model.Group, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	return m.group, nil
}

func (m *MockGroupForRecipe) validateGroup(group model.Group) error {
	if m.GetByIDError != nil {
		return m.GetByIDError
	}

	_ = group
	return nil
}

func TestNewRecipeService(t *testing.T) {
	mockRepo := &MockRecipeRepository{}

	service := NewRecipeService(
		mockRepo,
		&MockIngredientForRecipe{},
		&MockCategoryForRecipe{},
		&MockGroupForRecipe{},
	)

	if service == nil {
		t.Fatal("NewRecipeService() returned nil")
	}
}

func TestGetByItemID(t *testing.T) {
	tests := []struct {
		name     string
		itemID   int64
		expected []model.Recipe
		err      error
	}{
		{
			name:   "Multiple recipes",
			itemID: 10,
			expected: []model.Recipe{
				{ID: 1, Name: "Grilled Chicken", ImageURL: new("/static/recipes/chicken.jpg")},
				{ID: 3, Name: "Tomato Soup", ImageURL: new("/static/recipes/soup.jpg")},
			},
		},
		{
			name:     "No recipe",
			itemID:   10,
			expected: []model.Recipe{},
		},
		{
			name:     "repository error",
			itemID:   10,
			expected: nil,
			err:      customErrors.NewInternalError("failed to fetch recipes", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRecipeRepository{
				recipes:      tt.expected,
				getByItemErr: tt.err,
			}

			service := NewRecipeService(
				mockRepo,
				&MockIngredientForRecipe{},
				&MockCategoryForRecipe{},
				&MockGroupForRecipe{},
			)
			actual, err := service.GetByItemID(tt.itemID)

			if tt.err != nil {
				if !utils.CompareErrors(err, tt.err) {
					t.Fatalf("GetByItemID() error = %v, want %v", err, tt.err)
				}
				if actual != nil {
					t.Fatalf("GetByItemID() expected nil recipes on error, got %v", actual)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetByItemID() unexpected error = %v", err)
			}

			if len(actual) != len(tt.expected) {
				t.Fatalf("GetByItemID() returned %d recipes, expected %d", len(actual), len(tt.expected))
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("GetByItemID() recipes mismatch: got %v, want %v", actual, tt.expected)
			}
		})
	}
}
