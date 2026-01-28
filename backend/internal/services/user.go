package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, errors.New("Failed to fetch users: " + err.Error())
	}

	return users, nil
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		// Créer erreur custom
		return nil, fmt.Errorf("Failed to fetch user by ID %v: %v", id, err.Error())
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user by username %v: %v", username, err.Error())
	}

	return user, nil
}

func (s *UserService) CreateUser(user *models.User) (int64, error) {
	user.CreatedAt = time.Now()

	// Chercher si le nom d'utilisateur existe déjà
	// s'il existe, retourner une erreur (d'ailleurs créer des erreurs custom)
	// sinon c'est ok, on peut créer l'utilisateur

	// Hasher le password avant de le stocker en base de données

	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("Failed to create user %v: %v", user.Username, err.Error())
	}

	return id, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	// Checker si l'utilisateur existe à partir de l'ID

	err := s.repo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("Failed to update user %v: %v", user.ID, err.Error())
	}

	return nil
}

// Créer route pour changer le mot de passe (vérifier l'ancien mot de passe avant)

func (s *UserService) DeleteUser(id int64) error {
	// Checker si l'utilisateur existe à partir de l'ID

	err := s.repo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("Failed to delete user %v: %v", id, err.Error())
	}

	return nil
}
