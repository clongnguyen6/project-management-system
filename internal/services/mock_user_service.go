package services

import (
	"context"
	"example/project-management-system/internal/models"
)

// MockUserService is a mock implementation of UserService for testing.
type MockUserService struct {
	CreateUserFunc    func(ctx context.Context, user *models.User) error
	GetUserByIDFunc   func(ctx context.Context, id uint) (*models.User, error)
	GetAllUsersFunc   func(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	DeleteUserFunc    func(ctx context.Context, id uint) error
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) error {
	return m.CreateUserFunc(ctx, user)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return m.GetUserByIDFunc(ctx, id)
}

func (m *MockUserService) GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return m.GetAllUsersFunc(ctx, page, pageSize)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uint) error {
	return m.DeleteUserFunc(ctx, id)
}