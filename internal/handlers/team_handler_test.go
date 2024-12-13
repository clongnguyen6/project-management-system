package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"example/project-management-system/internal/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Services and Repositories
type MockTeamService struct {
	mock.Mock
}

func (m *MockTeamService) CreateTeam(ctx context.Context, team *models.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamService) GetTeamByID(ctx context.Context, id uint) (*models.Team, error) {
	args := m.Called(ctx, id)
	if team, ok := args.Get(0).(*models.Team); ok {
		return team, args.Error(1)
	}
	return nil, args.Error(1)

}

func (m *MockTeamService) GetPaginatedTeams(ctx context.Context, page, pageSize int) ([]models.Team, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if teams, ok := args.Get(0).([]models.Team); ok {
		return teams, int64(len(teams)), args.Error(1)
	}
	return args.Get(0).([]models.Team), args.Get(1).(int64), args.Error(2)
}

func (m *MockTeamService) UpdateTeam(ctx context.Context, team *models.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamService) DeleteTeam(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTeam(t *testing.T) {
	mockService := new(MockTeamService)
	handler := NewTeamHandler(mockService)

	t.Run("Successful Team Creation", func(t *testing.T) {
		team := &models.Team{
			Name:       "Test Team",
			Description: "Test Description",
			ProjectID:   1,
		}

		mockService.On("CreateTeam", mock.Anything, team).Return(nil)

		jsonTeam, _ := json.Marshal(team)
		req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(jsonTeam))
		w := httptest.NewRecorder()

		handler.CreateTeam(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		fmt.Println(w.Body)
		assert.Contains(t, w.Body.String(), `{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"name":"Test Team","description":"Test Description","project_id":1,"project":{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"name":"","description":"","start_date":"0001-01-01T00:00:00Z","end_date":"0001-01-01T00:00:00Z","status":"","user_ids":null,"users":null,"tasks":null,"teams":null},"users":null}`)
	})
}


func TestGetTeamByID(t *testing.T) {
	mockService := new(MockTeamService)
	handler := NewTeamHandler(mockService)

	t.Run("Successful Team Retrieval", func(t *testing.T) {
		team := &models.Team{
			Name:       "Test Team",
			Description: "Test Description",
			ProjectID:   1,
		}

		mockService.On("GetTeamByID", mock.Anything, uint(1)).Return(team, nil)

		req := httptest.NewRequest(http.MethodGet, "/teams/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.GetTeamByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllTeams(t *testing.T) {
	mockService := new(MockTeamService)
	handler := NewTeamHandler(mockService)

	t.Run("GetAllUsers", func(t *testing.T) {
		teams := []models.Team{
			{ID: 1, Name: "Team 1", Description: "Team 1 description"},
			{ID: 2, Name: "Team 2", Description: "Team 2 description"},
			{ID: 3, Name: "Team 3", Description: "Team 3 description"},
		}

		mockService.On("GetPaginatedTeams", mock.Anything, 1, 10).Return(teams, nil)

		req := httptest.NewRequest(http.MethodGet, "/teams", nil)
		w := httptest.NewRecorder()

		handler.GetPaginatedTeams(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}


func TestUpdateTeam(t *testing.T) {
	mockService := new(MockTeamService)
	handler := NewTeamHandler(mockService)

	t.Run("Successful Team Update", func(t *testing.T) {
		team := &models.Team{
			Name:       "Updated Team",
			Description: "Test Description",
			ProjectID:   1,
		}

		mockService.On("UpdateTeam", mock.Anything, team).Return(nil)

		jsonTeam, _ := json.Marshal(team)
		req := httptest.NewRequest(http.MethodPut, "/teams", bytes.NewBuffer(jsonTeam))
		w := httptest.NewRecorder()

		handler.UpdateTeam(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteTeam(t *testing.T) {
	mockService := new(MockTeamService)
	handler := NewTeamHandler(mockService)

	t.Run("Successful Team Deletion", func(t *testing.T) {
		mockService.On("DeleteTeam", mock.Anything, uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/teams/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.DeleteTeam(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}