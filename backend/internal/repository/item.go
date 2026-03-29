package repository

import (
	"database/sql"
	"errors"
	"strconv"

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

// GetAllByGroupID fetches all items by group ID, ordered by a specified column.
func (r *ItemRepository) GetAllByGroupID(groupID int64, sort string) ([]model.Item, error) {
	items := []model.Item{}

	rows, err := r.db.Query(`
	SELECT i.*, ic.id, ic.name
	FROM items i
	LEFT JOIN item_categories ic ON i.item_category_id = ic.id
	WHERE i.group_id = ?
	ORDER BY i.?`, groupID, sort)

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
			&item.UnitType,
			&item.GroupID,
			&item.ItemCategory.ID,
			&item.ItemCategory.Name,
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

// GetByID retrieves an item from the database by its ID.
func (r *ItemRepository) GetByID(id int64) (*model.Item, error) {
	item, err := r.fetchItem("id", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewNotFoundError("Item", strconv.FormatInt(id, 10), err)
		}
		return nil, customErrors.NewInternalError("Failed to fetch item", err)
	}

	return item, nil
}

// GetByName retrieves an item from the database by its name.
func (r *ItemRepository) GetByName(name string) (*model.Item, error) {
	item, err := r.fetchItem("name", name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewNotFoundError("Item", name, err)
		}
		return nil, customErrors.NewInternalError("Failed to fetch item", err)
	}

	return item, nil
}

// Create inserts a new item into the database and returns the inserted ID.
func (r *ItemRepository) Create(item *model.Item) (int64, error) {
	result, err := r.db.Exec(`
	INSERT INTO items (name, description, average_market_price, unit_type, item_category_id, group_id)
	VALUES (?, ?, ?, ?, ?, ?)`,
		item.Name,
		item.Description,
		item.AverageMarketPrice,
		item.UnitType,
		item.ItemCategory.ID,
		item.GroupID,
	)

	if err != nil {
		return 0, customErrors.NewInternalError("Failed to create item", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, customErrors.NewInternalError("Failed to retrieve last insert ID", err)
	}

	return id, nil
}

// Update modifies an existing item in the database.
func (r *ItemRepository) Update(item *model.Item) error {
	result, err := r.db.Exec(`
	UPDATE items
	SET name = ?, description = ?, average_market_price = ?, unit_type = ?, item_category_id = ?
	WHERE id = ?`,
		item.Name,
		item.Description,
		item.AverageMarketPrice,
		item.UnitType,
		item.ItemCategory.ID,
		item.ID,
	)

	if err != nil {
		return customErrors.NewInternalError("Failed to update item", err)
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalError("Failed to retrieve updated item", err)
	}

	// If no rows were updated, it means the item was not found.
	if updatedRow == 0 {
		return customErrors.NewNotFoundError("Item", strconv.FormatInt(item.ID, 10), err)
	}

	return nil
}

// Delete removes an item from the database by its ID.
func (r *ItemRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return customErrors.NewInternalError("Failed to delete item", err)
	}

	deletedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalError("Failed to retrieve deleted item", err)
	}

	// If no rows were deleted, it means the item was not found.
	if deletedRow == 0 {
		return customErrors.NewNotFoundError("Item", strconv.FormatInt(id, 10), err)
	}

	return nil
}

// fetchItem is a helper method to retrieve an item based on a specific column and value.
func (r *ItemRepository) fetchItem(column string, value any) (*model.Item, error) {
	item := &model.Item{}

	query := `
	SELECT i.*, ic.id, ic.name
	FROM items i
	LEFT JOIN item_categories ic ON i.item_category_id = ic.id
	WHERE i.` + column + ` = ?`

	row := r.db.QueryRow(query, value)

	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.AverageMarketPrice,
		&item.UnitType,
		&item.GroupID,
		&item.ItemCategory.ID,
		&item.ItemCategory.Name,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}
