package repository

import (
	"database/sql"
	"fmt"
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
	opt := utils.NewSelectFilteringOptions([]string{"recipes.id"}, []any{id}, "", false)
	recipes, err := r.fetchRecipes(opt)
	if err != nil {
		return nil, err
	}
	return &recipes[0], nil
}

func (r *RecipeRepository) fetchRecipes(opt *utils.SelectFilteringOptions) ([]model.Recipe, error) {
	query := fmt.Sprintf(`SELECT
	recipes.*,
	recipe_categories.id, recipe_categories.name,
	ingredients.id, ingredients.quantity, ingredients.item_id, ingredients.unit_id
	FROM recipes
	JOIN recipes_categories_junction ON recipes_categories_junction.recipe_id = recipes.id
	JOIN recipe_categories ON recipe_categories.id = recipes_categories_junction.category_id
	JOIN ingredients ON ingredients.recipe_id = recipes.id
	%s;`, utils.MakeSelectFiltering(opt))

	rows, err := r.db.Query(query, opt.WhereValues...)
	if err != nil {
		return nil, customErrors.NewInternalError("failed to fetch recipes", err)
	}

	m := make(map[int64]*model.Recipe)
	seenCategories := make(map[int64]map[int64]bool)
	seenIngredients := make(map[int64]map[int64]bool)
	tmpRecipe := &model.Recipe{}
	tmpCategory := &model.RecipeCategory{}
	tmpIngredient := &model.Ingredient{}

	for rows.Next() {
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
			&tmpIngredient.ItemID,
			&tmpIngredient.UnitID,
		)

		if err != nil {
			return nil, customErrors.NewInternalError("failed to fetch recipes", err)
		}

		id := tmpRecipe.ID
		if _, exists := m[id]; !exists {
			m[id] = tmpRecipe
			seenCategories[id] = make(map[int64]bool)
			seenIngredients[id] = make(map[int64]bool)
		}

		if !seenCategories[id][tmpCategory.ID] {
			m[id].Categories = append(m[id].Categories, *tmpCategory)
			seenCategories[id][tmpCategory.ID] = true
		}

		if !seenIngredients[id][tmpIngredient.ID] {
			m[id].Ingredients = append(m[id].Ingredients, *tmpIngredient)
			seenIngredients[id][tmpIngredient.ID] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("failed to fetch recipes", err)
	}

	if len(m) == 0 {
		return []model.Recipe{}, customErrors.NewNotFoundError("recipe", strings.Join(opt.WhereColumns, ","), err)
	}

	ret := make([]model.Recipe, 0, len(m))
	for _, recipe := range m {
		ret = append(ret, *recipe)
	}

	return ret, nil
}
