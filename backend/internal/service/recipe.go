package service

import (
	"context"
	"errors"
	"time"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type RecipeServiceInterface interface {
	GetByItemID(itemID int64) ([]model.Recipe, error)
	Create(ctx context.Context, recipe model.Recipe) (int64, error)
}

type RecipeService struct {
	repo              repository.RecipeRepositoryInterface
	ingredientService IngredientService
	categoryService   RecipeCategoryService
	groupService      GroupService
}

func NewRecipeService(recipeRepo repository.RecipeRepositoryInterface,
	ingredientService IngredientService,
	categoryService RecipeCategoryService,
	groupService GroupService) RecipeServiceInterface {
	return &RecipeService{
		repo:              recipeRepo,
		ingredientService: ingredientService,
		categoryService:   categoryService,
		groupService:      groupService,
	}
}

func (s *RecipeService) GetByItemID(itemID int64) ([]model.Recipe, error) {
	recipes, err := s.repo.GetByItemID(itemID)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

// Create validates the recipe and its ingredients before creating it in the database.
func (s *RecipeService) Create(ctx context.Context, recipe model.Recipe) (int64, error) {
	if err := s.validateRecipe(recipe); err != nil {
		return 0, err
	}

	for _, ing := range recipe.Ingredients {
		if err := s.ingredientService.validateIngredient(ing); err != nil {
			return 0, err
		}
	}

	for _, cat := range recipe.Categories {
		if err := s.categoryService.validateRecipeCategory(cat); err != nil {
			return 0, err
		}
	}

	recipe.CreatedAt = time.Now().UTC()

	id, err := s.repo.Create(ctx, &recipe)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// validateRecipe verifies the validity of the simple fields and the existence of the group associated with the recipe.
func (s *RecipeService) validateRecipe(recipe model.Recipe) error {
	if recipe.Name == "" {
		return customErrors.NewValidationError("name", "recipe must have a name", nil)
	}

	if recipe.Servings <= 0 {
		return customErrors.NewValidationError("servings", "recipe must have servings greater than 0", nil)
	}

	if recipe.PreparationTimeMin != nil && *recipe.PreparationTimeMin < 0 {
		return customErrors.NewValidationError("preparation_time_min", "preparation time cannot be negative", nil)
	}

	if recipe.CookingTimeMin != nil && *recipe.CookingTimeMin < 0 {
		return customErrors.NewValidationError("cooking_time_min", "cooking time cannot be negative", nil)
	}

	if len(recipe.Ingredients) == 0 {
		return customErrors.NewValidationError("ingredients", "recipe must have at least one ingredient", nil)
	}

	if recipe.ID != 0 {
		if _, err := s.groupService.GetByID(recipe.GroupID); err != nil {
			if _, isNotFoundError := errors.AsType[*customErrors.NotFoundError](err); isNotFoundError {
				return customErrors.NewConflictError("Group", "group must exists", nil)
			}
			return err
		}
	}

	return nil
}
