package services

import (
	"context"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/repositories"
	"fmt"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTaskByID(ctx context.Context, id uint) (*models.Task, error)
	GetTasksByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id uint) error
}

type TaskServiceImplementation struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &TaskServiceImplementation{repo: repo}
}

func (s *TaskServiceImplementation) CreateTask(ctx context.Context, task *models.Task) error {
	// Example validation logic
	fmt.Println("==task.Title: ", task.Title)
	if task.Title == "" {
		fmt.Println("==ERR: ", fmt.Errorf("task title is required"))
		return fmt.Errorf("task title is required")
	}
	fmt.Println("===ProjectID: ", task.ProjectID)
	if task.ProjectID == 0 {
		return fmt.Errorf("task must be associated with a project")
	}
	return s.repo.CreateTask(ctx, task)
}

func (s *TaskServiceImplementation) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	fmt.Println("AAAAAAAA")
	task, err := s.repo.GetTaskByID(ctx, id)
	fmt.Println("BBBBBBB")
	if err != nil {
		fmt.Println("Service ERR: ", err)
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (s *TaskServiceImplementation) GetTasksByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error) {
	return s.repo.GetTaskByProject(ctx, projectID, page, pageSize)
}

func (s *TaskServiceImplementation) UpdateTask(ctx context.Context, task *models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("task title is required")
	}
	return s.repo.UpdateTask(ctx, task)
}

func (s *TaskServiceImplementation) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.DeleteTask(ctx, id)
}
