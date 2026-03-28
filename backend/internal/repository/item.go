package repository

import (
	"database/sql"

	"github.com/zouipo/yumsday/backend/internal/model"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

type ItemRepositoryInterface interface {
	// Define methods for item data operations here
}

type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository constructs a new ItemRepository using the provided database.
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

// GetAll fetches all items from the database.
func (r *ItemRepository) GetAll() ([]model.Item, error) {
	items := []model.Item{}

	rows, err := r.db.Query("SELECT * FROM item")

	if err != nil {
		return nil, customErrors.NewInternalError("Failed to fetch items", err)
	}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.AverageMarketPrice,
			&item.ItemCategory,
			&item.UnitType,
		)

		if err != nil {
			rows.Close()
			return nil, customErrors.NewInternalError("Failed to scan item", err)
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("Failed to iterate rows", err)
	}

	return items, nil
}

func (r *ItemRepository) GetByID(id int64) (*model.Item, error) {
	var item model.Item

	err := r.db.QueryRow("SELECT * FROM item WHERE id = ?", id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.AverageMarketPrice,
		&item.ItemCategory,
		&item.UnitType,
	)

	if err != nil {
		return nil, customErrors.NewInternalError("Failed to fetch item", err)
	}

	return &item, nil
}
