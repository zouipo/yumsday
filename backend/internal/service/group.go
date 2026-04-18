package service

import (
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/repository"
)

type GroupService interface {
	GetByID(id int64) (*model.Group, error)
}

type groupService struct {
	repo repository.GroupRepositoryInterface
}

func NewGroupService(repo repository.GroupRepositoryInterface) GroupService {
	return &groupService{
		repo: repo,
	}
}

func (s *groupService) GetByID(id int64) (*model.Group, error) {
	group, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return group, nil
}
