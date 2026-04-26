package repository

import (
	"database/sql"
	"log/slog"

	"github.com/zouipo/yumsday/backend/internal/model"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

type ItemRepositoryInterface interface {
	GetByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error)
	GetByID(id int64) (*model.Item, error)
	GetByName(name string) ([]model.Item, error)
	Create(item *model.Item) (int64, error)
	Update(item *model.Item) error
	Delete(id int64) error
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
func (r *ItemRepository) GetByGroupID(groupID int64, sort string, desc bool) ([]model.Item, error) {
	clauses := "WHERE items.group_id = ? ORDER by " + sort

	if desc {
		clauses += " DESC"
	}

	items, err := r.fetchItems(clauses, groupID, sort)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetByID retrieves an item from the database by its ID.
func (r *ItemRepository) GetByID(id int64) (*model.Item, error) {
	items, err := r.fetchItems("WHERE items.id = ?", id)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, customErrors.NewNotFoundError("Item", "id", nil)
	}

	return &items[0], nil
}

// GetByName retrieves an item from the database by its name.
func (r *ItemRepository) GetByName(name string) ([]model.Item, error) {
	items, err := r.fetchItems("WHERE items.name LIKE concat('%', ?, '%') ORDER BY items.name", name)
	if err != nil {
		return nil, err
	}

	return items, nil
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
		return 0, customErrors.NewInternalError("failed to create item", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, customErrors.NewInternalError("failed to retrieve last insert ID", err)
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
		return customErrors.NewInternalError("failed to update item", err)
	}

	updatedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalError("failed to retrieve updated item", err)
	}

	// If no rows were updated, it means the item was not found.
	if updatedRow == 0 {
		return customErrors.NewNotFoundError("Item", "id", err)
	}

	return nil
}

// Delete removes an item from the database by its ID.
func (r *ItemRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return customErrors.NewInternalError("failed to delete item", err)
	}

	deletedRow, err := result.RowsAffected()
	if err != nil {
		return customErrors.NewInternalError("failed to retrieve deleted item", err)
	}

	// If no rows were deleted, it means the item was not found.
	if deletedRow == 0 {
		return customErrors.NewNotFoundError("Item", "id", err)
	}

	return nil
}

// fetchItems is a helper method to retrieve multiple items based on filtering options.
func (r *ItemRepository) fetchItems(clauses string, values ...any) ([]model.Item, error) {
	query := `SELECT 
	items.*, item_categories.name 
	FROM items
	LEFT JOIN item_categories ON items.item_category_id = item_categories.id ` + clauses

	slog.Debug("fetching items", "query", query)

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, customErrors.NewInternalError("failed to fetch items", err)
	}

	items := []model.Item{}

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
			return nil, customErrors.NewInternalError("failed to fetch items", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("failed to iterate rows", err)
	}

	return items, nil
}
