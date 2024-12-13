package services

import (
	"context"
	"example/project-management-system/internal/models"
)

type MockProjectService struct {
	CreateProjectFunc       func(ctx context.Context, project *models.Project) error
	GetProjectByIDFunc      func(ctx context.Context, id uint) (*models.Project, error)
	GetPaginatedProjectsFunc      func(ctx context.Context, page, pageSize int) ([]models.Project, int64, error)
	UpdateProjectFunc       func(ctx context.Context, project *models.Project) error
	DeleteProjectFunc       func(ctx context.Context, id uint) error
	GetTasksByProjectIDFunc func(ctx context.Context, projectID uint) ([]models.Task, error)
}

func (m *MockProjectService) CreateProject(ctx context.Context, project *models.Project) error {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectService) GetProjectByID(ctx context.Context, id uint) (*models.Project, error) {
	if m.GetProjectByIDFunc != nil {
		return m.GetProjectByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectService) GetPaginatedProjects(ctx context.Context, page, pageSize int) ([]models.Project, int64, error) {
	if m.GetPaginatedProjectsFunc != nil {
		return m.GetPaginatedProjectsFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockProjectService) UpdateProject(ctx context.Context, project *models.Project) error {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectService) DeleteProject(ctx context.Context, id uint) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectService) GetTaskByProjectID(ctx context.Context, projectID uint) ([]models.Task, error) {
	if m.GetTasksByProjectIDFunc != nil {
		return m.GetTasksByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}
