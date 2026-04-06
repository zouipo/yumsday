package repository

import (
	"database/sql"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	testUnit = map[int64]model.Unit{
		1:  {ID: 1, Name: "Kilogram"},
		2:  {ID: 2, Name: "Gram"},
		7:  {ID: 7, Name: "Teaspoon"},
		8:  {ID: 8, Name: "Piece"},
		11: {ID: 11, Name: "Undefined"},
	}

	testRecipes = []model.Recipe{
		{
			ID:                 1,
			Name:               "Grilled Chicken",
			Description:        new("Simple grilled chicken breast with herbs"),
			ImageURL:           new("/static/recipes/chicken.jpg"),
			PreparationTimeMin: new(10),
			CookingTimeMin:     new(20),
			Servings:           new(4),
			Instructions:       new("Season and grill until cooked through"),
			CreatedAt:          time.Unix(0, 0).UTC(),
			Public:             true,
			GroupID:            1,
			Categories: []model.RecipeCategory{
				{ID: 2, Name: "MAIN COURSE"},
			},
			Ingredients: []model.Ingredient{
				{ID: 1, Quantity: new(4.0), Unit: testUnit[8], Item: model.Item{ID: 7, Name: "Chicken Breast"}},
				{ID: 2, Quantity: new(2.0), Unit: testUnit[2], Item: model.Item{ID: 10, Name: "Garlic"}},
				{ID: 3, Quantity: new(0.5), Unit: testUnit[7], Item: model.Item{ID: 3, Name: "Salt"}},
			},
		},
		{
			ID:                 2,
			Name:               "Chocolate Chip Cookies",
			Description:        new("Classic homemade chocolate chip cookies"),
			ImageURL:           new("/static/recipes/cookies.jpg"),
			OriginalLink:       new("https://example.com/cookies"),
			PreparationTimeMin: new(15),
			CookingTimeMin:     new(12),
			Servings:           new(24),
			Instructions:       new("Mix ingredients and bake at 350F"),
			CreatedAt:          time.Unix(0, 0).UTC(),
			Public:             true,
			Comment:            new("Family favorite!"),
			GroupID:            1,
			Categories: []model.RecipeCategory{
				{ID: 1, Name: "DESSERT"},
			},
			Ingredients: []model.Ingredient{
				{ID: 4, Quantity: new(2.0), Unit: testUnit[1], Item: model.Item{ID: 1, Name: "Flour"}},
				{ID: 5, Quantity: new(1.0), Unit: testUnit[1], Item: model.Item{ID: 2, Name: "Sugar"}},
				{ID: 6, Quantity: new(0.5), Unit: testUnit[1], Item: model.Item{ID: 6, Name: "Butter"}},
				{ID: 7, Quantity: new(2.0), Unit: testUnit[8], Item: model.Item{ID: 4, Name: "Eggs"}},
			},
		},
		{
			ID:                 3,
			Name:               "Tomato Soup",
			Description:        new("Creamy tomato soup"),
			ImageURL:           new("/static/recipes/soup.jpg"),
			OriginalLink:       new("https://example.com/soup"),
			PreparationTimeMin: new(10),
			CookingTimeMin:     new(30),
			Servings:           new(6),
			Instructions:       new("Cook tomatoes with onions and blend"),
			CreatedAt:          time.Unix(0, 0).UTC(),
			Public:             false,
			Comment:            new("Great for winter"),
			GroupID:            2,
			Categories: []model.RecipeCategory{
				{ID: 3, Name: "SOUP"},
				{ID: 4, Name: "VEGETARIAN"},
			},
			Ingredients: []model.Ingredient{
				{ID: 8, Quantity: new(6.0), Unit: testUnit[8], Item: model.Item{ID: 8, Name: "Tomatoes"}},
				{ID: 9, Quantity: new(1.0), Unit: testUnit[8], Item: model.Item{ID: 9, Name: "Onions"}},
				{ID: 10, Quantity: new(2.0), Unit: testUnit[2], Item: model.Item{ID: 10, Name: "Garlic"}},
				{ID: 11, Quantity: new(1.0), Unit: testUnit[7], Item: model.Item{ID: 3, Name: "Salt"}},
			},
		},
		{
			ID:        4,
			Name:      "Quick Salad",
			CreatedAt: time.Unix(0, 0).UTC(),
			Public:    true,
			GroupID:   1,
			Categories: []model.RecipeCategory{
				{ID: 5, Name: "SALAD"},
				{ID: 7, Name: "VEGAN"},
			},
			Ingredients: []model.Ingredient{
				{ID: 12, Quantity: nil, Unit: testUnit[8], Item: model.Item{ID: 8, Name: "Tomatoes"}},
				{ID: 13, Quantity: new(1.0), Unit: testUnit[11], Item: model.Item{ID: 13, Name: "Olive Oil"}},
				{ID: 14, Quantity: nil, Unit: testUnit[11], Item: model.Item{ID: 12, Name: "Pepper"}},
			},
		},
	}
)

func setupRecipeTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?_foreign_keys=on")
	//db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Apply migrations using the migration package
	migrationsFS := os.DirFS("../../data/migrations")
	err = migration.Migrate(db, migrationsFS)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	testScript, _ := os.ReadFile("../../data/test.sql")
	_, err = db.Exec(string(testScript))
	if err != nil {
		t.Fatalf("failed to run test.sql: %v", err)
	}

	return db
}

func areRecipesEqual(r1 *model.Recipe, r2 *model.Recipe) bool {
	r1.Categories = utils.SortSliceByFieldName(r1.Categories, "ID", false)
	r2.Categories = utils.SortSliceByFieldName(r2.Categories, "ID", false)
	r1.Ingredients = utils.SortSliceByFieldName(r1.Ingredients, "ID", false)
	r2.Ingredients = utils.SortSliceByFieldName(r2.Ingredients, "ID", false)
	return reflect.DeepEqual(r1, r2)
}

func areRecipeSlicesEqual(s1 []model.Recipe, s2 []model.Recipe) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if !areRecipesEqual(&s1[i], &s2[i]) {
			return false
		}
	}

	return true
}

func TestGetByID(t *testing.T) {
	db := setupRecipeTestDB(t)
	defer db.Close()
	repo := NewRecipeRepository(db)

	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			"id 1",
			1,
			nil,
		},
		{
			"id 2",
			2,
			nil,
		},
		{
			"id 3",
			3,
			nil,
		},
		{
			"id 4",
			4,
			nil,
		},
		{
			"non existing id",
			-1,
			customErrors.NewNotFoundError("recipe", "recipes.id", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := repo.GetByID(tt.id)

			if tt.err != nil {
				if !utils.CompareErrors(err, tt.err) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("didn't expected error, got %v", err)
			}

			if !reflect.DeepEqual(*recipe, testRecipes[tt.id-1]) {
				t.Fatal("recipes should be equal")
			}
		})
	}
}

func TestGetByGroupID(t *testing.T) {
	db := setupRecipeTestDB(t)
	defer db.Close()
	repo := NewRecipeRepository(db)

	tests := []struct {
		name     string
		groupID  int64
		expected []model.Recipe
		err      error
	}{
		{
			name:    "group with one recipe",
			groupID: 2,
			expected: []model.Recipe{
				testRecipes[2],
			},
		},
		{
			name:    "group with multiple recipes",
			groupID: 1,
			expected: []model.Recipe{
				testRecipes[1],
				testRecipes[0],
				testRecipes[3],
			},
		},
		{
			name:     "group without recipe",
			groupID:  4,
			expected: []model.Recipe{},
		},
		{
			name:     "unknown group",
			groupID:  -1,
			expected: []model.Recipe{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.GetByGroupID(tt.groupID)

			if tt.err != nil {
				if !utils.CompareErrors(err, tt.err) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("didn't expected error, got %v", err)
			}

			if !areRecipeSlicesEqual(actual, tt.expected) {
				t.Fatal("recipes should be equal")
			}
		})
	}
}

func TestRecipeRepositoryCreate(t *testing.T) {
	db := setupRecipeTestDB(t)
	defer db.Close()
	repo := NewRecipeRepository(db)

	recipeID := new(int64(0))

	newRecipe := &model.Recipe{
		Name:               "test",
		Description:        new("description"),
		ImageURL:           new("http://example.com/test"),
		OriginalLink:       new("http://marmiton/test"),
		PreparationTimeMin: new(4),
		CookingTimeMin:     new(2),
		Servings:           new(1),
		Instructions:       new("[\"faire cuire !!\"]"),
		CreatedAt:          time.Now().UTC(),
		Public:             true,
		Comment:            new("comment !!"),
		GroupID:            1,
		Categories: []model.RecipeCategory{
			{
				ID:   1,
				Name: "DESSERT",
			},
			{
				ID:   2,
				Name: "MAIN COURSE",
			},
			{
				ID:   3,
				Name: "SOUP",
			},
		},
		Ingredients: []model.Ingredient{
			{
				Quantity: new(3.0),
				RecipeID: *recipeID,
				Item:     model.Item{ID: 1, Name: "Flour"},
				Unit:     model.Unit{ID: 1, Name: "Kilogram"},
			},
			{
				Quantity: new(3.0),
				RecipeID: *recipeID,
				Item:     model.Item{ID: 2, Name: "Sugar"},
				Unit:     model.Unit{ID: 2, Name: "Gram"},
			},
		},
	}

	id, err := repo.Create(newRecipe)
	if err != nil {
		t.Fatalf("expected no error, got '%s'", err)
	}

	*recipeID = id
	newRecipe.ID = id

	actual, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("expected no error, got '%s'", err)
	}

	// Or get the last two inserted IDs in the ingredients table instead
	for i := range actual.Ingredients {
		actual.Ingredients[i].ID = 0
	}

	if !reflect.DeepEqual(actual, newRecipe) {
		actualJson, _ := json.MarshalIndent(actual, "", "  ")
		expectedJson, _ := json.MarshalIndent(newRecipe, "", "  ")
		t.Fatalf("recipes should be equal: %s vs %s", actualJson, expectedJson)
	}
}
