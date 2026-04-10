package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

type RecipeRepositoryInterface interface {
	GetByID(id int64) (model.Recipe, error)
	GetByGroupID(groupID int64) ([]model.Recipe, error)
	Create(recipe *model.Recipe) (int64, error)
	Update(recipe *model.Recipe) error
	Delete(id int64) error
}

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
	}
}

func (r *RecipeRepository) GetByID(id int64) (*model.Recipe, error) {
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "recipes.id", Values: []any{id}},
		},
	}
	recipes, err := r.fetchRecipes(opt)
	if err != nil {
		return nil, err
	}
	return &recipes[0], nil
}

func (r *RecipeRepository) GetByGroupID(groupID int64) ([]model.Recipe, error) {
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "recipes.group_id", Values: []any{groupID}},
		},
		OrderBy: []utils.OrderByClause{
			{Column: "recipes.name"},
		},
	}
	recipes, err := r.fetchRecipes(opt)
	if err != nil {
		if _, isNotFoundError := errors.AsType[*customErrors.NotFoundError](err); isNotFoundError {
			return []model.Recipe{}, nil
		}
		return nil, err
	}
	return recipes, nil
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *model.Recipe, testHook func()) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO recipes(
			name,
			description,
			image_url,
			original_link,
			preparation_time_min,
			cooking_time_min,
			servings,
			instructions,
			created_at,
			public,
			comment,
			group_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
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
		return 0, customErrors.NewInternalError("Failed to create recipe", err)
	}

	recipe.ID, err = res.LastInsertId()
	if err != nil {
		return 0, customErrors.NewInternalError("Failed to retrieve recipe ID", err)
	}

	if err = r.createRecipeCategoryJunction(ctx, tx, recipe); err != nil {
		return 0, err
	}

	if err = r.createIngredients(ctx, tx, recipe); err != nil {
		return 0, err
	}

	tx.Commit()
	return recipe.ID, nil
}

func (r *RecipeRepository) fetchRecipes(opt *utils.SelectFilteringOptions) ([]model.Recipe, error) {
	query := fmt.Sprintf(`SELECT
	recipes.*,
	recipe_categories.id, recipe_categories.name,
	ingredients.id, ingredients.quantity,
	items.id, items.name,
	units.id, units.name
	FROM recipes
	LEFT JOIN recipes_categories_junction ON recipes_categories_junction.recipe_id = recipes.id
	LEFT JOIN recipe_categories ON recipe_categories.id = recipes_categories_junction.category_id
	LEFT JOIN ingredients ON ingredients.recipe_id = recipes.id
	LEFT JOIN items ON items.id = ingredients.item_id
	LEFT JOIN units ON units.id = ingredients.unit_id
	%s;`, utils.MakeSelectFiltering(opt))

	slog.Debug("fetching recipes", "query", query)

	rows, err := r.db.Query(query, opt.WhereValues()...)
	if err != nil {
		return nil, customErrors.NewInternalError("failed to fetch recipes", err)
	}

	ret := []model.Recipe{}

	// The query joins recipes_categories_junction and ingredients on recipes.id,
	// so the number of rows returned is a product
	// e.g. 3 categories * 5 ingredients = 15 rows returned.
	// We have to deduplicate all these rows lest we add duplicated stuff in our returned data.
	// This struct contains the data used to do this bookkeeping.
	type state struct {
		retIndex        int64
		seenCategories  map[int64]bool
		seenIngredients map[int64]bool
	}
	stateMap := make(map[int64]state)

	for rows.Next() {
		tmpRecipe := &model.Recipe{}
		tmpCategory := &model.RecipeCategory{}
		tmpIngredient := &model.Ingredient{}

		err := rows.Scan(
			&tmpRecipe.ID,
			&tmpRecipe.Name,
			&tmpRecipe.Description,
			&tmpRecipe.ImageURL,
			&tmpRecipe.OriginalLink,
			&tmpRecipe.PreparationTimeMin,
			&tmpRecipe.CookingTimeMin,
			&tmpRecipe.Servings,
			&tmpRecipe.Instructions,
			&tmpRecipe.CreatedAt,
			&tmpRecipe.Public,
			&tmpRecipe.Comment,
			&tmpRecipe.GroupID,
			&tmpCategory.ID,
			&tmpCategory.Name,
			&tmpIngredient.ID,
			&tmpIngredient.Quantity,
			&tmpIngredient.Item.ID,
			&tmpIngredient.Item.Name,
			&tmpIngredient.Unit.ID,
			&tmpIngredient.Unit.Name,
		)

		if err != nil {
			return nil, customErrors.NewInternalError("failed to fetch recipes", err)
		}

		id := tmpRecipe.ID
		if _, exists := stateMap[id]; !exists {
			ret = append(ret, *tmpRecipe)
			stateMap[id] = state{
				retIndex:        int64(len(ret) - 1),
				seenCategories:  make(map[int64]bool),
				seenIngredients: make(map[int64]bool),
			}
		}

		i := stateMap[id].retIndex

		if !stateMap[id].seenCategories[tmpCategory.ID] {
			ret[i].Categories = append(ret[i].Categories, *tmpCategory)
			stateMap[id].seenCategories[tmpCategory.ID] = true
		}

		if !stateMap[id].seenIngredients[tmpIngredient.ID] {
			ret[i].Ingredients = append(ret[i].Ingredients, *tmpIngredient)
			stateMap[id].seenIngredients[tmpIngredient.ID] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("failed to fetch recipes", err)
	}

	if len(ret) == 0 {
		return nil, customErrors.NewNotFoundError("recipe", strings.Join(opt.WhereColumns(), ","), nil)
	}

	return ret, nil
}

func (r *RecipeRepository) createRecipeCategoryJunction(ctx context.Context, tx *sql.Tx, recipe *model.Recipe) error {
	query := "INSERT INTO recipes_categories_junction(recipe_id, category_id) VALUES " +
		strings.Join(slices.Repeat([]string{"(?, ?)"}, len(recipe.Categories)), ", ")

	values := make([]any, 0, len(recipe.Categories)*2)
	for _, c := range recipe.Categories {
		values = append(values, recipe.ID, c.ID)
	}

	slog.Debug("inserting recipe category junctions", "query", query)

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return customErrors.NewInternalError("failed to insert recipe category junctions", err)
	}

	return nil
}

func (r *RecipeRepository) createIngredients(ctx context.Context, tx *sql.Tx, recipe *model.Recipe) error {
	query := "INSERT INTO ingredients(quantity, recipe_id, item_id, unit_id) VALUES " +
		strings.Join(slices.Repeat([]string{"(?, ?, ?, ?)"}, len(recipe.Ingredients)), ", ")

	values := make([]any, 0, len(recipe.Ingredients)*4)
	for _, i := range recipe.Ingredients {
		values = append(values, i.Quantity, recipe.ID, i.Item.ID, i.Unit.ID)
	}

	slog.Debug("inserting ingredients", "query", query)

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return customErrors.NewInternalError("failed to insert ingredients", err)
	}

	return nil
}
