package repository

import (
	"database/sql"
	"errors"
	"fmt"

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

func (r *RecipeRepository) fetchRecipes(column []string, value []any) ([]model.Recipe, error) {
	if len(column) != len(value) {
		panic("fetchRecipes: columns and values have different length")
	}

	query := `
	SELECT r.*, cat.id, cat.name, ing.id, ing.quantity, ing.item_id, ing.unit_id
	FROM recipes r
	JOIN recipes_categories_junction rcj ON rcj.recipe_id = r.id
	JOIN recipe_categories cat ON cat.id = rcj.category_id
	JOIN ingredients ing ON ing.recipe_id = r.id
	WHERE `

	for i := 0; i < len(column); i++ {
		query += fmt.Sprintf("%v = ? ", column)
		if i < len(column)-1 {
			query += "AND "
		}
	}
	query += ";"

	rows, err := r.db.Query(query, value...)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]*model.Recipe)
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
			if errors.Is(err, sql.ErrNoRows) {
				return []model.Recipe{}, nil
			}
			return nil, customErrors.NewInternalError("failed to fetch recipes", err)
		}

		id := tmpRecipe.ID
		if _, exists := m[id]; !exists {
			m[id] = tmpRecipe
		}

		m[id].Categories = append(m[id].Categories, *tmpCategory)
		m[id].Ingredients = append(m[id].Ingredients, *tmpIngredient)
	}

	ret := make([]model.Recipe, len(m))
	for _, recipe := range m {
		ret = append(ret, *recipe)
	}

	return ret, nil
}
