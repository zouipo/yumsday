package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

type ItemRepositoryInterface interface {
	GetAllByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error)
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
func (r *ItemRepository) GetByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error) {
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "items.group_id", Values: []any{groupID}},
		},
		OrderBy: []utils.OrderByClause{
			{Column: sort, Descending: descending},
		},
	}

	items, err := r.fetchItems(opt)
	if err != nil {
		// Returns an empty slice if no items were found for this group ID
		if _, isNotFoundError := errors.AsType[*customErrors.NotFoundError](err); isNotFoundError {
			return []model.Item{}, nil
		}
		return nil, err
	}

	return items, nil
}

// GetByID retrieves an item from the database by its ID.
func (r *ItemRepository) GetByID(id int64) (*model.Item, error) {
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "items.id", Values: []any{id}},
		},
	}

	items, err := r.fetchItems(opt)
	if err != nil {
		return nil, err
	}

	return &items[0], nil
}

// GetByName retrieves an item from the database by its name.
func (r *ItemRepository) GetByName(name string) ([]model.Item, error) {
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "items.name", Values: []any{name}},
		},
	}

	items, err := r.fetchItems(opt)
	if err != nil {
		// Returns an empty slice if no items were found for this group ID
		if _, isNotFoundError := errors.AsType[*customErrors.NotFoundError](err); isNotFoundError {
			return []model.Item{}, nil
		}
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
		return customErrors.NewNotFoundError("Item", "items.id", err)
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
		return customErrors.NewNotFoundError("Item", "items.id", err)
	}

	return nil
}

// fetchItems is a helper method to retrieve multiple items based on filtering options.
func (r *ItemRepository) fetchItems(opt *utils.SelectFilteringOptions) ([]model.Item, error) {
	query := fmt.Sprintf(`
	SELECT items.*, item_categories.name
	FROM items
	JOIN item_categories ON items.item_category_id = item_categories.id
	%s;`, utils.MakeSelectFiltering(opt))

	slog.Debug("fetching items", "query", query)

	rows, err := r.db.Query(query, opt.WhereValues()...)
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

	if len(items) == 0 {
		return nil, customErrors.NewNotFoundError("Item", strings.Join(opt.WhereColumns(), ","), err)
	}

	return items, nil
}
