package service

import (
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type ItemCategoryServiceInterface interface {
	GetByID(id int64) (*model.ItemCategory, error)
	GetByNameAndGroupID(name string, groupID int64) (*model.ItemCategory, error)
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
	return s.repo.GetByID(id)
}

func (s *ItemCategoryService) GetByNameAndGroupID(name string, groupID int64) (*model.ItemCategory, error) {
	return s.repo.GetByNameAndGroupID(name, groupID)
}
