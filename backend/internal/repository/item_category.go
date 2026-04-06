package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

type ItemCategoryRepositoryInterface interface {
	GetByID(id int64) (*model.ItemCategory, error)
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
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "item_categories.id", Values: []any{id}},
		},
	}

	itemCategories, err := r.fetchItemCategories(opt)
	if err != nil {
		return nil, err
	}

	return &itemCategories[0], nil
}

// fetchItemCategories is a helper method to retrieve multiple item categories based on filtering options.
func (r *ItemCategoryRepository) fetchItemCategories(opt *utils.SelectFilteringOptions) ([]model.ItemCategory, error) {
	query := fmt.Sprintf(`
	SELECT item_categories.*
	FROM item_categories
	%s;`, utils.MakeSelectFiltering(opt))

	slog.Debug("fetching item categories", "query", query)

	rows, err := r.db.Query(query, opt.WhereValues()...)
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

	if len(itemCategories) == 0 {
		return nil, customErrors.NewNotFoundError("ItemCategory", strings.Join(opt.WhereColumns(), ","), err)
	}

	return itemCategories, nil
}
