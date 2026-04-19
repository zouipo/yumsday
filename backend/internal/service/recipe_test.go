package service

import (
	"context"
	"reflect"
	"testing"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	invalidRecipeIngID   = -1
	invalidRecipeCatID   = -1
	invalidRecipeGroupID = -1

	pieceUnit = model.Unit{
		ID:   1,
		Name: "piece",
	}

	recipeGroup1 = model.Group{
		ID:   1,
		Name: "Family",
	}

	recipeGroup2 = model.Group{
		ID:   2,
		Name: "Friends",
	}

	validRecipeIngredients = []model.Ingredient{
		{
			ID:       1,
			Quantity: new(3.0),
			Item: model.Item{
				ID:      1,
				Name:    "Tomato",
				GroupID: recipeGroup1.ID,
			},
			Unit: pieceUnit,
		},
		{
			ID:       2,
			Quantity: new(2.5),
			Item: model.Item{
				ID:      2,
				Name:    "Onion",
				GroupID: recipeGroup1.ID,
			},
			Unit: pieceUnit,
		},
	}

	recipeIngredientGroup2 = model.Ingredient{
		ID:       3,
		Quantity: new(6.0),
		Item: model.Item{
			ID:      3,
			Name:    "Flour",
			GroupID: recipeGroup2.ID,
		},
		Unit: pieceUnit,
	}

	validRecipeCategories = []model.RecipeCategory{
		{
			ID:      1,
			Name:    "Dinner",
			GroupID: recipeGroup1.ID,
		},
		{
			ID:      2,
			Name:    "Vegetarian",
			GroupID: recipeGroup1.ID,
		},
	}

	categoryRecipe2 = model.RecipeCategory{
		ID:      3,
		Name:    "Dessert",
		GroupID: recipeGroup2.ID,
	}
)

type MockRecipeRepository struct {
	recipes      []model.Recipe
	getByItemErr error
	getByIDErr   error
	createErr    error
}

func NewMockRecipeRepository() *MockRecipeRepository {
	return &MockRecipeRepository{
		recipes: make([]model.Recipe, 0),
	}
}

func (m *MockRecipeRepository) GetByID(id int64) (*model.Recipe, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	for i := range m.recipes {
		if m.recipes[i].ID == id {
			return &m.recipes[i], nil
		}
	}

	return nil, customErrors.NewNotFoundError("Recipe", "id", nil)
}

func (m *MockRecipeRepository) GetByName(_ string, _ bool) ([]model.Recipe, error) {
	return []model.Recipe{}, nil
}

func (m *MockRecipeRepository) GetByGroupID(_ int64, _ bool) ([]model.Recipe, error) {
	return nil, nil
}

func (m *MockRecipeRepository) GetByItemID(itemID int64) ([]model.Recipe, error) {
	if m.getByItemErr != nil {
		return nil, m.getByItemErr
	}

	return m.recipes, nil
}

func (m *MockRecipeRepository) GetRecipeGroupID(_ int64) (int64, error) {
	return 0, nil
}

func (m *MockRecipeRepository) Create(_ context.Context, recipe *model.Recipe) (int64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}

	id := int64(len(m.recipes) + 1)
	recipe.ID = id
	m.recipes = append(m.recipes, *recipe)

	return id, nil
}

func (m *MockRecipeRepository) Update(_ context.Context, _ *model.Recipe) error {
	return nil
}

func (m *MockRecipeRepository) Delete(_ context.Context, _ int64) error {
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

func TestNewRecipeService(t *testing.T) {
	mockRepo := &MockRecipeRepository{}

	service := NewRecipeService(
		mockRepo,
		&MockIngredientForRecipe{},
		&MockCategoryForRecipe{},
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

func TestCreateRecipe(t *testing.T) {
	tests := []struct {
		name        string
		recipe      *model.Recipe
		expectedID  int64
		createErr   error
		ingErr      error
		catErr      error
		expectedErr error
	}{
		{
			name: "valid recipe with all fields",
			recipe: &model.Recipe{
				Name:               "Pasta Bolognese",
				Description:        new("Classic pasta with meat sauce"),
				ImageURL:           new("/static/recipes/pasta-bolognese.jpg"),
				OriginalLink:       new("https://example.com/pasta-bolognese"),
				PreparationTimeMin: new(30),
				CookingTimeMin:     new(30),
				Servings:           4,
				Instructions:       new("Cook pasta, simmer sauce, combine and serve."),
				Public:             true,
				Comment:            new("Family favorite"),
				GroupID:            recipeGroup1.ID,
				Ingredients:        validRecipeIngredients,
				Categories:         validRecipeCategories,
			},
			expectedID: 1,
		},
		{
			name: "valid recipe with nullable fields not provided",
			recipe: &model.Recipe{
				Name:        "Grilled Cheese",
				Servings:    2,
				GroupID:     recipeGroup1.ID,
				Ingredients: []model.Ingredient{validRecipeIngredients[0]},
				Categories:  []model.RecipeCategory{validRecipeCategories[0]},
			},
			expectedID: 1,
		},
		{
			name: "valid recipe with no category provided",
			recipe: &model.Recipe{
				Name:        "Scrambled Eggs",
				Servings:    1,
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
			},
			expectedID: 1,
		},
		{
			name: "Validation error empty name",
			recipe: &model.Recipe{
				Name:        "",
				Servings:    2,
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
			},
			expectedErr: customErrors.NewValidationError("name", "recipe must have a name", nil),
		},
		{
			name: "Validation error servings",
			recipe: &model.Recipe{
				Name:        "Soup",
				Servings:    -1,
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
			},
			expectedErr: customErrors.NewValidationError("servings", "recipe must have servings greater than 0", nil),
		},
		{
			name: "Validation error servings not provided",
			recipe: &model.Recipe{
				Name:        "Soup",
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
			},
			expectedErr: customErrors.NewValidationError("servings", "recipe must have servings greater than 0", nil),
		},
		{
			name: "Validation error preparation time less than 0",
			recipe: &model.Recipe{
				Name:               "Stew",
				Servings:           3,
				PreparationTimeMin: new(-1),
				GroupID:            recipeGroup1.ID,
				Ingredients:        validRecipeIngredients,
			},
			expectedErr: customErrors.NewValidationError("preparation_time_min", "preparation time cannot be negative", nil),
		},
		{
			name: "Validation error cooking time less than 0",
			recipe: &model.Recipe{
				Name:           "Stew",
				Servings:       3,
				CookingTimeMin: new(-1),
				GroupID:        recipeGroup1.ID,
				Ingredients:    validRecipeIngredients,
			},
			expectedErr: customErrors.NewValidationError("cooking_time_min", "cooking time cannot be negative", nil),
		},
		{
			name: "Validation error no ingredient provided",
			recipe: &model.Recipe{
				Name:     "Bread",
				Servings: 2,
				GroupID:  recipeGroup1.ID,
			},
			expectedErr: customErrors.NewValidationError("ingredients", "recipe must have at least one ingredient", nil),
		},
		{
			name: "Validation error invalid item for ingredient",
			recipe: &model.Recipe{
				Name:        "Bread",
				Servings:    2,
				GroupID:     recipeGroup1.ID,
				Ingredients: []model.Ingredient{recipeIngredientGroup2},
			},
			expectedErr: customErrors.NewConflictError("Ingredient", "item composing ingredient must belongs to the same group as the recipe", nil),
		},
		{
			name: "Invalid ingredient provided",
			recipe: &model.Recipe{
				Name:     "Burger",
				Servings: 2,
				GroupID:  recipeGroup1.ID,
				Ingredients: []model.Ingredient{
					{
						ID: int64(invalidRecipeIngID),
					},
				},
			},
			ingErr:      customErrors.NewValidationError("ingredient", "invalid ingredient", nil),
			expectedErr: customErrors.NewValidationError("ingredient", "invalid ingredient", nil),
		},
		{
			name: "Invalid category provided",
			recipe: &model.Recipe{
				Name:        "Burger",
				Servings:    2,
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
				Categories: []model.RecipeCategory{
					{
						ID: int64(invalidRecipeCatID),
					},
				},
			},
			catErr:      customErrors.NewValidationError("recipe category", "invalid recipe category", nil),
			expectedErr: customErrors.NewValidationError("recipe category", "invalid recipe category", nil),
		},
		{
			name: "Validation error invalid recipe category group",
			recipe: &model.Recipe{
				Name:        "Bread",
				Servings:    2,
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
				Categories:  []model.RecipeCategory{categoryRecipe2},
			},
			expectedErr: customErrors.NewConflictError("RecipeCategory", "recipe category must belongs to the same group as the recipe", nil),
		},
		{
			name: "Repository error",
			recipe: &model.Recipe{
				Name:        "Ramen",
				Servings:    2,
				GroupID:     recipeGroup1.ID,
				Ingredients: validRecipeIngredients,
				Categories:  validRecipeCategories,
			},
			createErr:   customErrors.NewInternalError("failed to create recipe", nil),
			expectedErr: customErrors.NewInternalError("failed to create recipe", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRecipeRepository{
				createErr: tt.createErr,
			}

			ingService := &MockIngredientForRecipe{
				validateIngError: tt.ingErr,
			}

			catService := &MockCategoryForRecipe{
				validateCatError: tt.catErr,
			}

			s := NewRecipeService(
				mockRepo,
				ingService,
				catService,
			)

			id, err := s.Create(context.Background(), tt.recipe)

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

			if id != 1 {
				t.Errorf("Create() returned ID 0, expected non-zero")
			}

			newRecipe, err := mockRepo.GetByID(1)
			if err != nil {
				t.Fatalf("GetByID() after Create() error: %v", err)
			}

			if newRecipe.CreatedAt.IsZero() {
				t.Errorf("CreatedAt attribute for new recipe shouldn't be zero value")
			}

			tt.recipe.CreatedAt = newRecipe.CreatedAt

			if !reflect.DeepEqual(newRecipe, tt.recipe) {
				t.Errorf("Recipes should be equal: expected %v, got %v", tt.recipe, newRecipe)
			}
		})
	}
}
