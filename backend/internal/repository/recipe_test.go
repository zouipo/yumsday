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
	testCategories = []model.RecipeCategory{
		{
			Name: "Dessert",
		},
		{
			Name: "Main course",
		},
		{
			Name: "Cake",
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
			Quantity: utils.Ptr(10.0),
			RecipeID: 2,
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
			GroupID:            1,
		},
		{
			Name:               "Olive Cake",
			Description:        utils.Ptr("Olive cake yummy"),
			ImageURL:           utils.Ptr("http://example.com/image2"),
			OriginalLink:       nil,
			PreparationTimeMin: utils.Ptr(60),
			CookingTimeMin:     utils.Ptr(30),
			Servings:           utils.Ptr(8),
			Instructions:       utils.Ptr("[\"Cake it damn it !!\"]"),
			CreatedAt:          time.Now().UTC(),
			Public:             false,
			Comment:            nil,
			GroupID:            1,
		},
	}

	recipeCategoryJunction = map[*model.Recipe][]*model.RecipeCategory{
		&testRecipes[0]: {
			&testCategories[0], &testCategories[2],
		},
		&testRecipes[1]: {
			&testCategories[1], &testCategories[2],
		},
	}
)

func setupRecipeTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
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

		// fill this recipe's ingredients list
		for _, ing := range testIngredients {
			if ing.RecipeID == testRecipes[i].ID {
				ing.RecipeID = 0
				testRecipes[i].Ingredients = append(recipe.Ingredients, ing)
			}
		}
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
			k.Categories = append(k.Categories, *cat)
		}
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
			"valid id",
			1,
			nil,
		},
		{
			"valid id 2",
			2,
			nil,
		},
		{
			"non existing id",
			-1,
			customErrors.NewNotFoundError("recipe", "id", nil),
		},
	}

	for _, tt := range tests {
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
	}
}
