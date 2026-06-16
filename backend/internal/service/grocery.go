package service

import "github.com/zouipo/yumsday/backend/internal/repository"

type GroceryServiceInterface interface {
	HasItem(id int64) (bool, error)
}

type GroceryService struct {
	repo repository.GroceryRepositoryInterface
}

func NewGroceryService(repo repository.GroceryRepositoryInterface) *GroceryService {
	return &GroceryService{
		repo: repo,
	}
}

func (s *GroceryService) HasItem(itemID int64) (bool, error) {
	return s.repo.HasItem(itemID)
}
