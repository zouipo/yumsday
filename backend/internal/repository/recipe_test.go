package repository

import (
	"database/sql"
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
	testRecipes = []model.Recipe{
		{
			ID:                 1,
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
				{ID: 1, Quantity: new(2.0), UnitID: 1, Item: model.Item{ID: 1, Name: "Flour"}},
				{ID: 2, Quantity: new(1.0), UnitID: 1, Item: model.Item{ID: 2, Name: "Sugar"}},
				{ID: 3, Quantity: new(0.5), UnitID: 1, Item: model.Item{ID: 6, Name: "Butter"}},
				{ID: 4, Quantity: new(2.0), UnitID: 8, Item: model.Item{ID: 4, Name: "Eggs"}},
			},
		},
		{
			ID:                 2,
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
				{ID: 5, Quantity: new(4.0), UnitID: 8, Item: model.Item{ID: 7, Name: "Chicken Breast"}},
				{ID: 6, Quantity: new(2.0), UnitID: 2, Item: model.Item{ID: 10, Name: "Garlic"}},
				{ID: 7, Quantity: new(0.5), UnitID: 7, Item: model.Item{ID: 3, Name: "Salt"}},
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
				{ID: 8, Quantity: new(6.0), UnitID: 8, Item: model.Item{ID: 8, Name: "Tomatoes"}},
				{ID: 9, Quantity: new(1.0), UnitID: 8, Item: model.Item{ID: 9, Name: "Onions"}},
				{ID: 10, Quantity: new(2.0), UnitID: 2, Item: model.Item{ID: 10, Name: "Garlic"}},
				{ID: 11, Quantity: new(1.0), UnitID: 7, Item: model.Item{ID: 3, Name: "Salt"}},
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
				{ID: 12, Quantity: nil, UnitID: 8, Item: model.Item{ID: 8, Name: "Tomatoes"}},
				{ID: 13, Quantity: new(1.0), UnitID: 11, Item: model.Item{ID: 13, Name: "Olive Oil"}},
				{ID: 14, Quantity: nil, UnitID: 11, Item: model.Item{ID: 12, Name: "Pepper"}},
			},
		},
	}
)

func setupRecipeTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?_foreign_keys=on")
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
