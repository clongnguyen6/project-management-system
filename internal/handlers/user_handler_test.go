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

func TestUserHandlers(t *testing.T) {
	mockUsers := []models.User{
		{BaseModel: models.BaseModel{ID: 1}, Username: "user1", Email: "user1@example.com"},
		{BaseModel: models.BaseModel{ID: 2}, Username: "user2", Email: "user2@example.com"},
	}
	mockUser := models.User{BaseModel: models.BaseModel{ID: 1}, Username: "testuser", Email: "test@example.com"}

	mockService := &services.MockUserService{
		CreateUserFunc: func(ctx context.Context, user *models.User) error {
			user.ID = 1 // Simulate the user being created with an ID
			return nil
		},
		GetUserByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
			if id == mockUser.ID {
				return &mockUser, nil
			}
			return nil, nil
		},
		GetAllUsersFunc: func(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
			return mockUsers, int64(len(mockUsers)), nil
		},
		DeleteUserFunc: func(ctx context.Context, id uint) error {
			return nil
		},
	}

	handler := NewUserHandler(mockService)

	t.Run("CreateUser", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)

		var createdUser models.User
		err := json.NewDecoder(res.Body).Decode(&createdUser)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", createdUser.Username)
		assert.Equal(t, "test@example.com", createdUser.Email)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		handler.GetAllUsers(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, float64(1), response["page"])
		assert.Equal(t, float64(2), response["total"])

		users := response["users"].([]interface{})
		assert.Equal(t, "user1", users[0].(map[string]interface{})["username"])
		assert.Equal(t, "user2", users[1].(map[string]interface{})["username"])
	})

	t.Run("GetUserByID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		req.SetPathValue("id", "1")

		w := httptest.NewRecorder()
		handler.GetUserByID(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var user models.User
		err := json.NewDecoder(res.Body).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		req.SetPathValue("id", "1")

		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var response map[string]string
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "user delete successfully", response["message"])
	})
}
