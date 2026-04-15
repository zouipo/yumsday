package service

import (
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type GroupServiceInterface interface {
	GetByID(id int64) (*model.Group, error)
}

type GroupService struct {
	repo repository.GroupRepositoryInterface
}

func NewGroupService(repo repository.GroupRepositoryInterface) *GroupService {
	return &GroupService{
		repo: repo,
	}
}

func (s *GroupService) GetByID(id int64) (*model.Group, error) {
	group, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return group, nil
}
