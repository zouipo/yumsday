package repository

import (
	"database/sql"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

type GroceryRepositoryInterface interface {
	HasItem(id int64) (bool, error)
}

type GroceryRepository struct {
	db *sql.DB
}

func NewGroceryRepository(db *sql.DB) *GroceryRepository {
	return &GroceryRepository{
		db: db,
	}
}

func (r *GroceryRepository) HasItem(itemID int64) (bool, error) {
	var exists bool
	query := `
	SELECT EXISTS(
	SELECT 1 FROM groceries
	WHERE item_id = ?)`

	err := r.db.QueryRow(query, itemID).Scan(&exists)
	if err != nil {
		return false, customErrors.NewInternalError("failed to get info from groceries", err)
	}

	return exists, nil
}
