package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectHandlers(t *testing.T) {
	mockProject := models.Project{
		BaseModel: models.BaseModel{ID: 1},
		Name:      "Test Project",
		Description: "Test Description",
	}

	mockProjects := []models.Project{
		{BaseModel: models.BaseModel{ID: 1}, Name: "Project Alpha", Description: "Description Alpha"},
		{BaseModel: models.BaseModel{ID: 2}, Name: "Project Beta", Description: "Description Beta"},
	}

	mockTasks := []models.Task{
		{BaseModel: models.BaseModel{ID: 1}, Title: "Task 1", Description: "Description for Task 1"},
		{BaseModel: models.BaseModel{ID: 2}, Title: "Task 2", Description: "Description for Task 2"},
	}


	mockService := &services.MockProjectService{
		CreateProjectFunc: func(ctx context.Context, project *models.Project) error {
			project.ID = 1 // Simulate project creation with an ID
			return nil
		},
		GetProjectByIDFunc: func(ctx context.Context, id uint) (*models.Project, error) {
			if id == mockProject.ID {
				return &mockProject, nil
			}
			return nil, nil
		},
		GetPaginatedProjectsFunc: func(ctx context.Context, page, pageSize int) ([]models.Project, int64, error) {
			return mockProjects, int64(len(mockProjects)), nil
		},
		DeleteProjectFunc: func(ctx context.Context, id uint) error {
			return nil
		},
		GetTasksByProjectIDFunc: func(ctx context.Context, projectID uint) ([]models.Task, error) {
			if projectID == 1 {
				return mockTasks, nil
			}
			return nil, nil
		},

	}

	handler := NewProjectHandler(mockService)

	t.Run("CreateProject", func(t *testing.T) {
		project := models.Project{
			Name:        "Test Project",
			Description: "A project for testing",
		}

		body, _ := json.Marshal(project)
		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateProject(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)

		var createdProject models.Project
		err := json.NewDecoder(res.Body).Decode(&createdProject)
		assert.NoError(t, err)
		assert.Equal(t, "Test Project", createdProject.Name)
	})


	t.Run("GetAllProjects", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects?page=1&pageSize=2", nil)
		w := httptest.NewRecorder()

		handler.GetAllProjects(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, float64(2), response["total"])
		projects := response["projects"].([]interface{})
		assert.Equal(t, "Project Alpha", projects[0].(map[string]interface{})["name"])
		assert.Equal(t, "Project Beta", projects[1].(map[string]interface{})["name"])
	})

	t.Run("GetProjectByID", func(t *testing.T) {
		mockProject := models.Project{BaseModel: models.BaseModel{ID: 1}, Name: "Test Project", Description: "Test Description"}

		mockService.GetProjectByIDFunc = func(ctx context.Context, id uint) (*models.Project, error) {
			if id == mockProject.ID {
				return &mockProject, nil
			}
			return nil, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
		req.SetPathValue("id", "1")

		w := httptest.NewRecorder()

		handler.GetProjectByID(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var project models.Project
		err := json.NewDecoder(res.Body).Decode(&project)
		assert.NoError(t, err)
		assert.Equal(t, "Test Project", project.Name)
	})

	t.Run("DeleteProject", func(t *testing.T) {
		mockService.DeleteProjectFunc = func(ctx context.Context, id uint) error {
			return nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/projects/1", nil)
		req.SetPathValue("id", "1")

		w := httptest.NewRecorder()

		handler.DeleteProject(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var response map[string]string
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "project deleted successfully", response["message"])
	})
}
