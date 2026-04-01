package repository

import (
	"database/sql"
	"errors"
	"strconv"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
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

func (r *ItemCategoryRepository) GetByID(id int64) (*model.ItemCategory, error) {
	category, err := r.fetchCategoryItem("id", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewNotFoundError("Item category", strconv.FormatInt(id, 10), err)
		}
		return nil, customErrors.NewInternalError("Failed to fetch item category", err)
	}

	return category, nil
}

func (r *ItemCategoryRepository) fetchCategoryItem(column string, value any) (*model.ItemCategory, error) {
	category := &model.ItemCategory{}

	query := `
	SELECT id, name
	FROM item_categories
	WHERE ` + column + ` = ?`

	row := r.db.QueryRow(query, value)

	err := row.Scan(
		&category.ID,
		&category.Name,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}
