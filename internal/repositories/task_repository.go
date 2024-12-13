package repositories

import (
	"context"
	"example/project-management-system/internal/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTaskByID(ctx context.Context, id uint) (*models.Task, error)
	GetTaskByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id uint) error
}

type TaskRepositoryImplementation struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImplementation{db: db}
}

func (r *TaskRepositoryImplementation) CreateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepositoryImplementation) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).
	Preload("Assignee").
	Preload("Project").
	First(&task, id).Error
	return &task, err
}

func (r *TaskRepositoryImplementation) GetTaskByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	// Count total records for the project
	if err := r.db.WithContext(ctx).Model(&models.Task{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated tasks
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Offset(offset).
		Limit(pageSize).
		Preload("User").
		Preload("Project").
		Find(&tasks).Error

	return tasks, total, err
}

func (r *TaskRepositoryImplementation) UpdateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *TaskRepositoryImplementation) DeleteTask(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}
