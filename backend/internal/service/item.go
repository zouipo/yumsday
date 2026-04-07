package service

import (
	"strings"

	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

type ItemServiceInterface interface {
	GetAllByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error)
	GetByID(id int64) (*model.Item, error)
	GetByName(name string) (*model.Item, error)
	Create(item *model.Item) (int64, error)
	Update(item *model.Item) error
	Delete(id int64) error
}

type ItemService struct {
	repo           repository.ItemRepositoryInterface
	recipeService  RecipeServiceInterface
	groceryService GroceryServiceInterface
}

// NewItemService creates a new ItemService using the provided ItemRepository.
func NewItemService(itemRepo repository.ItemRepositoryInterface,
	recipeService RecipeServiceInterface,
	groceryService GroceryServiceInterface) *ItemService {
	return &ItemService{
		repo:           itemRepo,
		recipeService:  recipeService,
		groceryService: groceryService,
	}
}

/*** READ OPERATIONS ***/
// GetAllByGroupID returns all items for a given group ID, sorted by the specified key and order.
func (s *ItemService) GetAllByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error) {
	sortKey, err := s.mapSortKey(strings.ToLower(sort))
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetAllByGroupID(groupID, sortKey, descending)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetByID returns the item identified by id or an error if not found.
func (s *ItemService) GetByID(id int64) (*model.Item, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetByName returns the item that matches the provided name or an error.
func (s *ItemService) GetByName(name string) ([]model.Item, error) {
	if name == "" {
		return nil, customErrors.NewNotFoundError("Item", name, nil)
	}

	items, err := s.repo.GetByName(name)
	if err != nil {
		return nil, err
	}

	return items, nil
}

/*** CREATE OPERATIONS ***/
func (s *ItemService) Create(item *model.Item) (int64, error) {

	id, err := s.repo.Create(item)
	if err != nil {
		return 0, err
	}

	return id, nil
}

/*** UPDATE OPERATIONS ***/
func (s *ItemService) Update(item *model.Item) error {
	err := s.repo.Update(item)
	if err != nil {
		return err
	}

	return nil
}

/*** DELETE OPERATIONS ***/
// Delete removes the item identified by id from the database.
// It checks for any dependencies in recipes and groceries before deletion.
func (s *ItemService) Delete(id int64) error {
	r, err := s.recipeService.GetByItemID(id)
	if err != nil {
		return err
	}

	if len(r) > 0 {
		return customErrors.NewConflictError("Item", "can't delete item used by recipes", nil)
	}

	b, err := s.groceryService.HasItem(id)
	if err != nil {
		return err
	}

	if b {
		return customErrors.NewConflictError("Item", "can't delete item used in groceries", nil)
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

/*** HELPER FUNCTIONS ***/
// mapSortKey maps the sort parameter to the corresponding database column.
func (s *ItemService) mapSortKey(param string) (string, error) {
	switch param {
	case "name", "":
		return "i.name", nil
	case "average_market_price":
		return "i.average_market_price", nil
	case "unit_type":
		return "i.unit_type", nil
	case "category":
		return "ic.name", nil
	default:
		return "", customErrors.NewInvalidParamsError(param, nil)
	}
}
