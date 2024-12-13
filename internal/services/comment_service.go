package services

import (
	"context"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/repositories"
	"fmt"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByID(ctx context.Context, id uint) (*models.Comment, error)
	GetCommentsByTask(ctx context.Context, taskID uint, page, pageSize int) ([]models.Comment, int64, error)
	DeleteComment(ctx context.Context, id uint) error
}

type CommentServiceImplementation struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) CommentService {
	return &CommentServiceImplementation{repo: repo}
}

func (s *CommentServiceImplementation) CreateComment(ctx context.Context, comment *models.Comment) error {
	// Validate input
	if comment.Content == "" {
		return fmt.Errorf("content is required")
	}
	if comment.TaskID == 0 {
		return fmt.Errorf("task ID is required")
	}
	return s.repo.CreateComment(ctx, comment)
}

func (s *CommentServiceImplementation) GetCommentByID(ctx context.Context, id uint) (*models.Comment, error) {
	comment, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}
	return comment, nil
}

func (s *CommentServiceImplementation) GetCommentsByTask(ctx context.Context, taskID uint, page, pageSize int) ([]models.Comment, int64, error) {
	return s.repo.GetCommentsByTask(ctx, taskID, page, pageSize)
}

func (s *CommentServiceImplementation) DeleteComment(ctx context.Context, id uint) error {
	return s.repo.DeleteComment(ctx, id)
}
