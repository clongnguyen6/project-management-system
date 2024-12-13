package repositories

import (
	"context"
	"example/project-management-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id uint) (*models.Project, error)
	GetPaginatedProjects(ctx context.Context, page, pageSize int) ([]models.Project, int64, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id uint) error
	GetTaskByProjectID(ctx context.Context, projectID uint) ([]models.Task, error)
}

type ProjectRepositoryImplementation struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &ProjectRepositoryImplementation{db: db}
}

func (r *ProjectRepositoryImplementation) CreateProject(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepositoryImplementation) GetProjectByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project

	if err := r.db.WithContext(ctx).
		Preload("Users").
		Preload("Tasks").
		Preload("Teams").
		First(&project, id).Error; err != nil {
		return nil, err
	}

	project.UserIDs = make([]uint, len(project.Users))
	for i, user := range project.Users {
		project.UserIDs[i] = user.ID
	}

	return &project, nil
}

func (r *ProjectRepositoryImplementation) GetPaginatedProjects(ctx context.Context, page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Project{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	// Fetch paginated records
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Preload("Users").
		Preload("Tasks").
		Preload("Teams").
		Offset(offset).
		Limit(pageSize).
		Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch projects: %w", err)
	}

	for i := range projects {
		projects[i].UserIDs = make([]uint, len(projects[i].Users))
		for j, user := range projects[i].Users {
			projects[i].UserIDs[j] = user.ID
		}
	}

	return projects, total, nil

}

func (r *ProjectRepositoryImplementation) UpdateProject(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepositoryImplementation) DeleteProject(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, id).Error
}

// AddUsersToProject add user to project
func (r *ProjectRepositoryImplementation) AddUsersToProject(ctx context.Context, projectID uint, userIDs []uint) error {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, projectID).Error; err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users, userIDs).Error; err != nil {
		return fmt.Errorf("error finding users: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&project).Association("Users").Append(users); err != nil {
		return fmt.Errorf("failed to add users to project: %w", err)
	}

	return nil
}

// RemoveUsersFromProject remove user from project
func (r *ProjectRepositoryImplementation) RemoveUsersFromProject(ctx context.Context, projectID uint, userIDs []uint) error {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, projectID).Error; err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users, userIDs).Error; err != nil {
		return fmt.Errorf("error finding users: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&project).Association("Users").Delete(users); err != nil {
		return fmt.Errorf("failed to remove users from project: %w", err)
	}

	return nil
}

func (r *ProjectRepositoryImplementation) SearchProjects(ctx context.Context, query string, page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	searchQuery := r.db.WithContext(ctx).
		Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")

	if err := searchQuery.Model(&models.Project{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	offset := (page - 1) * pageSize

	if err := searchQuery.
		Preload("Users").
		Preload("Tasks").
		Preload("Teams").
		Offset(offset).
		Limit(pageSize).
		Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search projects: %w", err)
	}

	for i := range projects {
		projects[i].UserIDs = make([]uint, len(projects[i].Users))
		for j, user := range projects[i].Users {
			projects[i].UserIDs[j] = user.ID
		}
	}

	return projects, total, nil
}

func (r *ProjectRepositoryImplementation) GetTaskByProjectID(ctx context.Context, projectID uint) ([]models.Task, error) {
	var tasks []models.Task

	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&tasks).Error

	return tasks, err
}