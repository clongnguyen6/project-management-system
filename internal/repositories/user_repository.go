package repositories

import (
	"context"
	"example/project-management-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository defines the methods for interacting with the user data.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	DeleteUser(ctx context.Context, id uint) error
}

// UserRepositoryImplementation is an implementation of the UserRepository using Gorm.
type UserRepositoryImplementation struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImplementation{db: db}
}

func (r *UserRepositoryImplementation) CreateUser(ctx context.Context, user *models.User) error {
    // Begin transaction for additional safety
    tx := r.db.WithContext(ctx).Begin()

	if len(user.ProjectIDs) > 0 {
		if err := tx.Find(&user.Projects, user.ProjectIDs).Error; err != nil {
            return fmt.Errorf("error loading projects: %w", err)
        }
	}

    // Attempt to create user
    if err := tx.Create(user).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create user: %w", err)
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}


func (r *UserRepositoryImplementation) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Projects").First(&user, id).Error; err != nil {
		return nil, err
	}
	// Create ProjectIDs from Projects
    user.ProjectIDs = make([]uint, len(user.Projects))
    for i, project := range user.Projects {
        user.ProjectIDs[i] = project.ID
    }

    return &user, nil
}

func (r *UserRepositoryImplementation) GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total users
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Fetch users with pagination
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Projects").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}

	for i := range users {
		users[i].ProjectIDs = make([]uint, len(users[i].Projects))
		for j, project := range users[i].Projects {
			users[i].ProjectIDs[j] = project.ID
		}
	}

	return users, total, nil
}

func (r *UserRepositoryImplementation) DeleteUser(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
