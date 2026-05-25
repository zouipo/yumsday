package service

import (
	"errors"

	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
)

type ItemServiceInterface interface {
	GetAllByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error)
	GetByID(id int64) (*model.Item, error)
	GetByName(name string, descending bool) ([]model.Item, error)
	GetRecipes(id int64) ([]model.Recipe, error)
	Create(item *model.Item) (int64, error)
	Update(item *model.Item) error
	Delete(id int64) error
}

type ItemService struct {
	repo                repository.ItemRepositoryInterface
	recipeService       RecipeServiceInterface
	groceryService      GroceryServiceInterface
	groupService        GroupServiceInterface
	itemCategoryService ItemCategoryServiceInterface
}

// NewItemService creates a new ItemService using the provided ItemRepository.
func NewItemService(itemRepo repository.ItemRepositoryInterface,
	recipeService RecipeServiceInterface,
	groceryService GroceryServiceInterface,
	groupService GroupServiceInterface,
	itemCategoryService ItemCategoryServiceInterface) *ItemService {
	return &ItemService{
		repo:                itemRepo,
		recipeService:       recipeService,
		groceryService:      groceryService,
		groupService:        groupService,
		itemCategoryService: itemCategoryService,
	}
}

/*** READ OPERATIONS ***/
// GetAllByGroupID returns all items for a given group ID, sorted by the specified key and order.
func (s *ItemService) GetByGroupID(groupID int64, sort string, descending bool) ([]model.Item, error) {
	return s.repo.GetByGroupID(groupID, sort, descending)
}

// GetByID returns the item identified by id or an error if not found.
func (s *ItemService) GetByID(id int64) (*model.Item, error) {
	return s.repo.GetByID(id)
}

// GetByName returns the item that matches the provided name or an error.
func (s *ItemService) GetByName(name string, descending bool) ([]model.Item, error) {
	if name == "" {
		return nil, customErrors.NewNotFoundError("Item", "name", nil)
	}

	return s.repo.GetByName(name, descending)
}

// GetRecipes returns the recipes in which the item is used.
func (s *ItemService) GetRecipes(id int64) ([]model.Recipe, error) {
	return s.recipeService.GetRecipeLiteByItemID(id)
}

/*** CREATE OPERATIONS ***/
// Create adds a new item to the database.
func (s *ItemService) Create(item *model.Item) (int64, error) {
	// if no item category is provided, assign the default one (uncategorized)
	if item.ItemCategory.ID == 0 {
		uncategorized, err := s.itemCategoryService.GetByNameAndGroupID("Uncategorized", item.GroupID)
		if err != nil {
			return 0, err
		}
		item.ItemCategory = *uncategorized
	}

	if err := s.validateItem(item); err != nil {
		return 0, err
	}

	return s.repo.Create(item)
}

/*** UPDATE OPERATIONS ***/
// Update modifies the item identified by id with the provided item data.
func (s *ItemService) Update(item *model.Item) error {
	currentItem, err := s.repo.GetByID(item.ID)
	if err != nil {
		return err
	}

	// GroupID can't be updated
	item.GroupID = currentItem.GroupID

	// If no item category is provided, assign the default one (uncategorized)
	if item.ItemCategory.ID == 0 {
		uncategorized, err := s.itemCategoryService.GetByNameAndGroupID("Uncategorized", item.GroupID)
		if err != nil {
			return err
		}
		item.ItemCategory = *uncategorized

		if err := checkSimpleFields(item); err != nil {
			return err
		}

	} else {
		if err := s.validateItem(item); err != nil {
			return err
		}
	}

	return s.repo.Update(item)
}

/*** DELETE OPERATIONS ***/
// Delete removes the item identified by id from the database.
// It checks for any dependencies in recipes and groceries before deletion.
func (s *ItemService) Delete(id int64) error {
	r, err := s.recipeService.GetRecipeLiteByItemID(id)
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

	return s.repo.Delete(id)
}

/*** HELPER FUNCTIONS ***/
// validateItem checks the validity of the item fields and ensures that related entities exist and are consistent.
func (s *ItemService) validateItem(item *model.Item) error {
	err := checkSimpleFields(item)

	if err != nil {
		return err
	}

	// If the item is new (create route)
	if item.ID == 0 {
		if _, err = s.groupService.GetByID(item.GroupID); err != nil {
			if _, isNotFoundError := errors.AsType[*customErrors.NotFoundError](err); isNotFoundError {
				return customErrors.NewConflictError("Group", "group must exists", nil)
			}
			return err
		}
	}

	itemCategory, err := s.itemCategoryService.GetByID(item.ItemCategory.ID)

	// Checks if the item category exists
	if err != nil {
		if _, isNotFoundError := errors.AsType[*customErrors.NotFoundError](err); isNotFoundError {
			return customErrors.NewConflictError("ItemCategory", "item category must exists", nil)
		}
		return err
	}

	if itemCategory.GroupID != item.GroupID {
		return customErrors.NewConflictError("ItemCategory", "item category must belongs to the same group as the item", nil)
	}

	return nil
}

// checkSimpleFields validates the basic fields of the item that don't require database access.
func checkSimpleFields(item *model.Item) error {
	e := customErrors.NewInvalidParamsError([]string{}, nil).(*customErrors.InvalidParamsError)

	if item.Name == "" {
		e.AddInvalidField("name")
	}

	// Caught when the field unit_type is omitted in the JSON body (set to the zero value, an empty string)
	if item.UnitType.String() == "" {
		e.AddInvalidField("unit type")
	}

	if len(e.Fields) > 0 {
		return e
	}

	return nil
}
