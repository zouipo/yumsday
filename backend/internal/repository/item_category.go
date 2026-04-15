package repository

import (
	"database/sql"
	"log/slog"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
)

type ItemCategoryRepositoryInterface interface {
	GetByID(id int64) (*model.ItemCategory, error)
	GetByNameAndGroupID(name string, groupID int64) (*model.ItemCategory, error)
}

type ItemCategoryRepository struct {
	db *sql.DB
}

func NewItemCategoryRepository(db *sql.DB) *ItemCategoryRepository {
	return &ItemCategoryRepository{
		db: db,
	}
}

// GetByID retrieves an item category from the database by its ID.
func (r *ItemCategoryRepository) GetByID(id int64) (*model.ItemCategory, error) {
	itemCategories, err := r.fetchItemCategories("WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if len(itemCategories) == 0 {
		return nil, customErrors.NewNotFoundError("item_categories", "id", nil)
	}

	return &itemCategories[0], nil
}

// GetByNameAndGroupID retrieves an item category from the database by its name and group ID.
func (r *ItemCategoryRepository) GetByNameAndGroupID(name string, groupID int64, descending bool) ([]model.ItemCategory, error) {
	clauses := "WHERE name LIKE concat('%', ?, '%') AND group_id = ? ORDER BY name"

	if descending {
		clauses += " DESC"
	}

	itemCategories, err := r.fetchItemCategories(clauses, name, groupID)
	if err != nil {
		return nil, err
	}

	return itemCategories, nil
}

// fetchItemCategories is a helper method to retrieve multiple item categories based on filtering options.
func (r *ItemCategoryRepository) fetchItemCategories(clauses string, values ...any) ([]model.ItemCategory, error) {
	query := `SELECT 
	item_categories.*
	FROM item_categories ` + clauses

	slog.Debug("fetching item categories", "query", query)

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, customErrors.NewInternalError("failed to fetch item categories", err)
	}

	itemCategories := []model.ItemCategory{}

	for rows.Next() {
		var itemCategory model.ItemCategory
		err := rows.Scan(
			&itemCategory.ID,
			&itemCategory.Name,
			&itemCategory.GroupID,
		)

		if err != nil {
			return nil, customErrors.NewInternalError("failed to fetch item categories", err)
		}

		itemCategories = append(itemCategories, itemCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("failed to iterate rows", err)
	}

	return itemCategories, nil
}
