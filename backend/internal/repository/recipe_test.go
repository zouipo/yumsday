package repository

import (
	"database/sql"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/zouipo/yumsday/backend/internal/migration"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

var (
	testCategories = []model.RecipeCategory{
		{
			Name:    "Dessert",
			GroupID: 1,
		},
		{
			Name:    "Main course",
			GroupID: 1,
		},
		{
			Name:    "Cake",
			GroupID: 1,
		},
	}

	testIngredients = []model.Ingredient{
		{
			Quantity: utils.Ptr(3.0),
			RecipeID: 1,
			ItemID:   2,
			UnitID:   3,
		},
		{
			Quantity: utils.Ptr(3.0),
			RecipeID: 4,
			ItemID:   5,
			UnitID:   6,
		},
	}

	testRecipes = []model.Recipe{
		{
			Name:               "Cheesecake",
			Description:        utils.Ptr("Le cheesecake super bon"),
			ImageURL:           utils.Ptr("http://example.com/image1"),
			OriginalLink:       utils.Ptr("http://recipe.com/cheesecake"),
			PreparationTimeMin: utils.Ptr(30),
			CookingTimeMin:     utils.Ptr(55),
			Servings:           utils.Ptr(8),
			Instructions:       utils.Ptr("[\"Prepare cheesecake\", \"Cook it\"]"),
			CreatedAt:          time.Now().UTC(),
			Public:             true,
			Comment:            utils.Ptr("best cheesecake !"),
			GroupID:            0,
			Categories:         []model.RecipeCategory{},
			Ingredients:        []model.Ingredient{},
		},
	}

	recipeCategoryJunction = map[*model.Recipe][]*model.RecipeCategory{
		&testRecipes[0]: {
			&testCategories[0], &testCategories[2],
		},
	}
)

func setupRecipeTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Apply migrations using the migration package
	migrationsFS := os.DirFS("../../data/migrations")
	err = migration.Migrate(db, migrationsFS)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	for i, category := range testCategories {
		res, err := db.Exec(
			`INSERT INTO recipe_categories(name, group_id) VALUES(?, ?);`,
			category.Name,
			category.GroupID,
		)
		if err != nil {
			t.Fatalf("failed to insert test recipe category '%s': %v", category.Name, err)
		}
		testCategories[i].ID, _ = res.LastInsertId()
	}

	for i, ingredient := range testIngredients {
		res, err := db.Exec(
			`INSERT INTO ingredients(quantity, recipe_id, item_id, unit_id) VALUES(?, ?, ?, ?);`,
			ingredient.Quantity,
			ingredient.RecipeID,
			ingredient.ItemID,
			ingredient.UnitID,
		)
		if err != nil {
			t.Fatalf("failed to insert test ingredient '%d': %v", ingredient.ID, err)
		}
		testIngredients[i].ID, _ = res.LastInsertId()
	}

	for i, recipe := range testRecipes {
		res, err := db.Exec(
			`INSERT INTO recipes(name, description, image_url, original_link, preparation_time_min, cooking_time_min, servings, instructions, created_at, public, comment, group_id)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
			recipe.Name,
			recipe.Description,
			recipe.ImageURL,
			recipe.OriginalLink,
			recipe.PreparationTimeMin,
			recipe.CookingTimeMin,
			recipe.Servings,
			recipe.Instructions,
			recipe.CreatedAt,
			recipe.Public,
			recipe.Comment,
			recipe.GroupID,
		)
		if err != nil {
			t.Fatalf("failed to insert test recipe '%s': %v", recipe.Name, err)
		}
		testRecipes[i].ID, _ = res.LastInsertId()
	}

	for k, v := range recipeCategoryJunction {
		for _, cat := range v {
			_, err := db.Exec(
				`INSERT INTO recipes_categories_junction(recipe_id, category_id)
				VALUES(?, ?);`,
				k.ID,
				cat.ID,
			)
			if err != nil {
				t.Fatalf(
					"failed to insert junction between recipe '%d' and category '%d': %v",
					k.ID,
					cat.ID,
					err,
				)
			}
		}
	}

	return db
}

func TestGetByID(t *testing.T) {
	db := setupRecipeTestDB(t)
	defer db.Close()

	repo := NewRecipeRepository(db)
	recipe, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("didn't expected error, got %v", err)
	}
	recipe.CreatedAt = testRecipes[0].CreatedAt
	if !reflect.DeepEqual(recipe, testRecipes[0]) {
		t.Fatal("recipes should be equal")
	}
}
