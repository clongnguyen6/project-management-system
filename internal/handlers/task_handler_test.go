package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"example/project-management-system/internal/models"
)

// Mock Services and Repositories
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskService) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	args := m.Called(ctx, id)
	// return args.Get(0).(*models.Task), args.Error(1)
	if task, ok := args.Get(0).(*models.Task); ok {
		return task, args.Error(1)
	}
	return nil, args.Error(1)

}

func (m *MockTaskService) GetTasksByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	return args.Get(0).([]models.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Task Handler Tests
func TestCreateTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	t.Run("Successful Task Creation", func(t *testing.T) {
		task := &models.Task{
			Title:       "Test Task",
			Description: "Test Description",
			ProjectID:   1,
			AssignedTo:  1,
		}

		mockService.On("CreateTask", mock.Anything, task).Return(nil)

		jsonTask, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonTask))
		w := httptest.NewRecorder()

		handler.CreateTask(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"title":"Test Task","description":"Test Description","project_id":1,"project":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"name":"","description":"","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","status":"","user_ids":null,"users":null,"tasks":null,"teams":null},"assigned_to":1,"assignee":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"username":"","email":"","first_name":"","last_name":"","project_ids":null,"projects":null,"role":""}}`)
	})

	t.Run("Invalid Task Creation", func(t *testing.T) {
		invalidTask := `{"title": 1.0}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(invalidTask))
		w := httptest.NewRecorder()
		handler.CreateTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTaskByID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	t.Run("Successful Task Retrieval", func(t *testing.T) {
		task := &models.Task{
			BaseModel: models.BaseModel{ID: 1},
			Title:     "Test Task",
		}

		mockService.On("GetTaskByID", mock.Anything, uint(1)).Return(task, nil)

		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.GetTaskByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Task Not Found", func(t *testing.T) {
		mockService.On("GetTaskByID", mock.Anything, uint(999)).Return(nil, errors.New("task not found"))

		req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		handler.GetTaskByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetTasksByProject(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	t.Run("Successful Tasks Retrieval", func(t *testing.T) {
		tasks := []models.Task{
			{
				BaseModel:    models.BaseModel{ID: 1},
				Title:        "Task 1",
				ProjectID:    1,
			},
			{
				BaseModel:    models.BaseModel{ID: 2},
				Title:        "Task 2",
				ProjectID:    1,
			},
		}

		mockService.On("GetTasksByProject", mock.Anything, uint(1), 1, 10).Return(tasks, int64(2), nil)

		req := httptest.NewRequest(http.MethodGet, "/projects/1/tasks?page=1&page_size=10", nil)
		req.SetPathValue("project_id", "1")
		w := httptest.NewRecorder()

		handler.GetTasksByProject(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	t.Run("Successful Task Update", func(t *testing.T) {
		task := &models.Task{
			BaseModel:    models.BaseModel{ID: 1},
			Title:        "Updated Task",
			ProjectID:    1,
		}

		mockService.On("UpdateTask", mock.Anything, task).Return(nil)

		jsonTask, _ := json.Marshal(task)
		req := httptest.NewRequest(http.MethodPut, "/tasks", bytes.NewBuffer(jsonTask))
		w := httptest.NewRecorder()

		handler.UpdateTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	t.Run("Successful Task Deletion", func(t *testing.T) {
		mockService.On("DeleteTask", mock.Anything, uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.DeleteTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

// // Service Layer Tests
// func TestTaskService_CreateTask(t *testing.T) {		
// 	mockService := new(MockTaskService)
// 	service := NewTaskService(mockService)

// 	t.Run("Invalid Task Creation", func(t *testing.T) {


// 		task := &models.Task{Title: ""}
// 		err := service.CreateTask(context.Background(), task)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "task title is required")
// 	})

// 	t.Run("Valid Task Creation", func(t *testing.T) {
// 		mockRepo := repositories.NewMockTaskRepository()
// 		service := services.NewTaskService(mockRepo)

// 		task := &models.Task{
// 			Title:       "Test Task",
// 			ProjectID:   1,
// 			Description: "Test Description",
// 		}

// 		mockRepo.On("CreateTask", mock.Anything, task).Return(nil)
// 		err := service.CreateTask(context.Background(), task)

// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})
// }
