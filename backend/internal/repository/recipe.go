package repository

import (
	"database/sql"
	"fmt"
	"strings"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
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
	recipes, err := r.fetchRecipes([]string{"id"}, []any{id})
	if err != nil {
		return nil, err
	}
	return &recipes[0], nil
}

func (r *RecipeRepository) fetchRecipes(columns []string, values []any) ([]model.Recipe, error) {
	if len(columns) != len(values) {
		panic("fetchRecipes: columns and values have different length")
	}

	query := `SELECT r.*, cat.id, cat.name, ing.id, ing.quantity, ing.item_id, ing.unit_id
	FROM recipes r
	JOIN recipes_categories_junction rcj ON rcj.recipe_id = r.id
	JOIN recipe_categories cat ON cat.id = rcj.category_id
	JOIN ingredients ing ON ing.recipe_id = r.id
	WHERE `

	for i := 0; i < len(columns); i++ {
		query += fmt.Sprintf("r.%v = ? ", columns[i])
		if i < len(columns)-1 {
			query += "AND "
		}
	}
	query += ";"

	rows, err := r.db.Query(query, values...)
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
		return []model.Recipe{}, customErrors.NewNotFoundError("recipe", strings.Join(columns, ","), err)
	}

	ret := make([]model.Recipe, 0, len(m))
	for _, recipe := range m {
		ret = append(ret, *recipe)
	}

	return ret, nil
}
