package services

import (
	"context"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/repositories"
	"fmt"
)

type ProjectService interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id uint) (*models.Project, error)
	GetPaginatedProjects(ctx context.Context, page, pageSize int) ([]models.Project, int64, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id uint) error
	GetTaskByProjectID(ctx context.Context, projectID uint) ([]models.Task, error)
}

type ProjectServiceImplementation struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &ProjectServiceImplementation{repo: repo}
}

func (s *ProjectServiceImplementation) CreateProject(ctx context.Context, project *models.Project) error {
	// Example: Additional validation logic
	if len(project.Name) < 3 {
		return fmt.Errorf("project name must be at least 3 characters")
	}
	return s.repo.CreateProject(ctx, project)
}

func (s *ProjectServiceImplementation) GetProjectByID(ctx context.Context, id uint) (*models.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}
	return project, nil
}

func (s *ProjectServiceImplementation) GetPaginatedProjects(ctx context.Context, page, pageSize int) ([]models.Project, int64, error) {
	return s.repo.GetPaginatedProjects(ctx, page, pageSize)
}

func (s *ProjectServiceImplementation) UpdateProject(ctx context.Context, project *models.Project) error {
	// Example: Additional validation logic
	if len(project.Name) < 3 {
		return fmt.Errorf("project name must be at least 3 characters")
	}
	return s.repo.UpdateProject(ctx, project)
}

func (s *ProjectServiceImplementation) DeleteProject(ctx context.Context, id uint) error {
	return s.repo.DeleteProject(ctx, id)
}

func (s *ProjectServiceImplementation) GetTaskByProjectID(ctx context.Context, projectID uint) ([]models.Task, error) {
	return s.repo.GetTaskByProjectID(ctx, projectID)
}