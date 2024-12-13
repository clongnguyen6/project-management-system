package repositories

import (
	"example/project-management-system/internal/models"

	"gorm.io/gorm"
)

type UserProjectRepository interface {
	AddUserToProject(userID uint, projectID uint) error
}

type UserProjectRepositoryImplementation struct {
	db *gorm.DB
}

func NewUserProjectRepository(db *gorm.DB) UserProjectRepository {
	return &UserProjectRepositoryImplementation{db: db}
}

// AddUserToProject adds a user to a project.
func (repo *UserProjectRepositoryImplementation) AddUserToProject(userID uint, projectID uint) error {
	user := &models.User{}
	project := &models.Project{}

	if err := repo.db.First(user, userID).Error; err != nil {
		return err
	}
	if err := repo.db.First(project, projectID).Error; err != nil {
		return err
	}

	return repo.db.Model(user).Association("Projects").Append(project)
}