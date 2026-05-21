package service

import (
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type RecipeServiceInterface interface {
	GetRecipeLiteByItemID(itemID int64) ([]model.Recipe, error)
}

type RecipeService struct {
	repo repository.RecipeRepositoryInterface
}

func NewRecipeService(recipeRepo repository.RecipeRepositoryInterface) *RecipeService {
	return &RecipeService{
		repo: recipeRepo,
	}
}

func (s *RecipeService) GetRecipeLiteByItemID(itemID int64) ([]model.Recipe, error) {
	recipes, err := s.repo.GetRecipeLiteByItemID(itemID)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}
