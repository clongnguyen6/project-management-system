package repositories

import (
	"context"
	"example/project-management-system/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, id uint) (*models.Comment, error)
	GetCommentsByTask(ctx context.Context, taskID uint, page, pageSize int) ([]models.Comment, int64, error)
	DeleteComment(ctx context.Context, id uint) error
}

type CommentRepositoryImplementation struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &CommentRepositoryImplementation{db: db}
}

func (r *CommentRepositoryImplementation) CreateComment(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepositoryImplementation) GetCommentByID(ctx context.Context, id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Task").
		First(&comment, id).Error
	return &comment, err
}

func (r *CommentRepositoryImplementation) GetCommentsByTask(ctx context.Context, taskID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	// Count total records for the task
	if err := r.db.WithContext(ctx).Model(&models.Comment{}).Where("task_id = ?", taskID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated comments
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Offset(offset).
		Limit(pageSize).
		Preload("User").
		Preload("Task").
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepositoryImplementation) DeleteComment(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}
