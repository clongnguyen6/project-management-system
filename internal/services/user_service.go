package services

import (
	"context"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/repositories"
)

// UserService defines the methods for performing business operations on Users.
type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	DeleteUser(ctx context.Context, id uint) error
}

// UserServiceImplementation is an implementation of the UserService.
type UserServiceImplementation struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImplementation{userRepo: userRepo}
}

func (s *UserServiceImplementation) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserServiceImplementation) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *UserServiceImplementation) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserServiceImplementation) GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.GetAllUsers(ctx, page, pageSize)
}

