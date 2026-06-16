package service

import (
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type RecipeServiceInterface interface {
	GetByItemID(itemID int64, descending bool) ([]model.Recipe, error)
}

type RecipeService struct {
	repo repository.RecipeRepositoryInterface
}

func NewRecipeService(recipeRepo repository.RecipeRepositoryInterface) *RecipeService {
	return &RecipeService{
		repo: recipeRepo,
	}
}

func (s *RecipeService) GetByItemID(itemID int64, descending bool) ([]model.Recipe, error) {
	recipes, err := s.repo.GetByItemID(itemID, descending)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}
