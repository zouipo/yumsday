package service

import (
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type ItemCategoryServiceInterface interface {
	GetByID(id int64) (*model.ItemCategory, error)
}

type ItemCategoryService struct {
	repo repository.ItemCategoryRepositoryInterface
}

func NewItemCategoryService(repo repository.ItemCategoryRepositoryInterface) *ItemCategoryService {
	return &ItemCategoryService{
		repo: repo,
	}
}

func (s *ItemCategoryService) GetByID(id int64) (*model.ItemCategory, error) {
	ic, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return ic, nil
}
