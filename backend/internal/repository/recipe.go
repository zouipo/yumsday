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
	tx, _ := r.db.BeginTx(ctx, nil)
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

	if err = r.updateRecipesCategoriesJunction(ctx, tx, recipe); err != nil {
		return 0, err
	}

	if err = r.updateIngredients(ctx, tx, recipe); err != nil {
		return 0, err
	}

	tx.Commit()
	return recipe.ID, nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *model.Recipe) error {
	tx, _ := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	tx.ExecContext(
		ctx,
		`UPDATE recipes
		SET
			name = ?,
			description = ?,
			image_url = ?,
			original_link = ?,
			preparation_time_min = ?,
			cooking_time_min = ?,
			servings = ?,
			instructions = ?,
			created_at = ?,
			public = ?,
			comment = ?
		WHERE id = ?`,
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
		recipe.ID,
	)

	if err := r.updateRecipesCategoriesJunction(ctx, tx, recipe); err != nil {
		return err
	}
	if err := r.updateIngredients(ctx, tx, recipe); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (r *RecipeRepository) Delete(ctx context.Context, id int64) error {
	tx, _ := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	deleteFunc := func(ctx context.Context, tx *sql.Tx, tableName string, columnName string, recipeID int64) error {
		res, err := tx.ExecContext(
			ctx,
			fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, columnName),
			recipeID,
		)
		if err != nil {
			return customErrors.NewInternalError(
				fmt.Sprintf("failed to delete %s by %s for recipe %d", tableName, columnName, recipeID),
				err,
			)
		}

		deletedRows, err := res.RowsAffected()
		if err != nil {
			return customErrors.NewInternalError(
				fmt.Sprintf("failed to check if recipe's %s were deleted", tableName),
				err,
			)
		}
		if deletedRows == 0 {
			return customErrors.NewNotFoundError(tableName, columnName, err)
		}

		return nil
	}

	if err := deleteFunc(ctx, tx, "recipes_categories_junction", "recipe_id", id); err != nil {
		if _, ok := errors.AsType[*customErrors.NotFoundError](err); ok {
			slog.WarnContext(ctx, "removing a recipe that has no recipes_categories_junction", "id", id)
		} else {
			return err
		}
	}
	if err := deleteFunc(ctx, tx, "recipes_dishes_junction", "recipe_id", id); err != nil {
		if _, ok := errors.AsType[*customErrors.NotFoundError](err); !ok {
			return err
		}
	}
	if err := deleteFunc(ctx, tx, "ingredients", "recipe_id", id); err != nil {
		if _, ok := errors.AsType[*customErrors.NotFoundError](err); ok {
			slog.WarnContext(ctx, "removing a recipe that has no ingredients", "id", id)
		} else {
			return err
		}
	}
	if err := deleteFunc(ctx, tx, "recipes", "id", id); err != nil {
		return err
	}

	tx.Commit()
	return nil
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

func (r *RecipeRepository) updateRecipesCategoriesJunction(ctx context.Context, tx *sql.Tx, recipe *model.Recipe) error {
	query := "INSERT INTO recipes_categories_junction (recipe_id, category_id) VALUES " +
		strings.Join(slices.Repeat([]string{"(?, ?)"}, len(recipe.Categories)), ", ") + " " +
		`ON CONFLICT(recipe_id,category_id) DO NOTHING`

	values := make([]any, 0, len(recipe.Categories)*2)
	for _, c := range recipe.Categories {
		values = append(values, recipe.ID, c.ID)
	}

	slog.Debug("update recipe category junctions", "query", query)

	if _, err := tx.ExecContext(ctx, query, values...); err != nil {
		return customErrors.NewInternalError("failed to update recipe category junctions", err)
	}

	query = `DELETE FROM recipes_categories_junction
			WHERE recipe_id = ? AND (recipe_id, category_id) NOT IN (` +
		strings.Join(slices.Repeat([]string{"(?, ?)"}, len(recipe.Categories)), ", ") + ")"

	values = slices.Concat([]any{recipe.ID}, values)

	slog.Debug("deleting obsolete recipes_categories_junction", "query", query)

	if _, err := tx.ExecContext(ctx, query, values...); err != nil {
		return customErrors.NewInternalError("failed to delete obsolete recipes_categories_junction", err)
	}

	return nil
}

func (r *RecipeRepository) updateIngredients(ctx context.Context, tx *sql.Tx, recipe *model.Recipe) error {
	upsertValues := make([]any, 0, len(recipe.Ingredients)*5)
	for _, ing := range recipe.Ingredients {
		var id any = ing.ID
		// If ingredient's ID is 0, we set it to nil which is interpreted as NULL by the db
		// that way it generates a new ID for it.
		if ing.ID == 0 {
			id = nil
		}
		upsertValues = append(upsertValues, id, ing.Quantity, recipe.ID, ing.Item.ID, ing.Unit.ID)
	}

	query := "INSERT INTO ingredients (id, quantity, recipe_id, item_id, unit_id) VALUES " +
		strings.Join(slices.Repeat([]string{"(?, ?, ?, ?, ?)"}, len(recipe.Ingredients)), ", ") + " " +
		`ON CONFLICT(id) DO UPDATE SET
			quantity = EXCLUDED.quantity,
			item_id = EXCLUDED.item_id,
			unit_id = EXCLUDED.unit_id
			RETURNING id`

	slog.Debug("update ingredients", "query", query)

	rows, err := tx.QueryContext(ctx, query, upsertValues...)
	if err != nil {
		return customErrors.NewInternalError("failed to update ingredients", err)
	}

	deleteValues := make([]any, 0, len(recipe.Ingredients)+1)
	deleteValues = append(deleteValues, recipe.ID)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return customErrors.NewInternalError("failed to update ingredients", err)
		}
		deleteValues = append(deleteValues, id)
	}

	query = `DELETE FROM ingredients 
			WHERE recipe_id = ? AND id NOT IN (` +
		strings.Join(slices.Repeat([]string{"?"}, len(recipe.Ingredients)), ", ") + ")"

	slog.Debug("deleting obsolete ingredients", "query", query)

	// deleteValues contains the ID of the recipe and its updated ingredients
	if _, err := tx.ExecContext(ctx, query, deleteValues...); err != nil {
		return customErrors.NewInternalError("failed to delete obsolete ingredients", err)
	}

	return nil
}
