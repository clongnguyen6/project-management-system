package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"example/project-management-system/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentService mocks the CommentService for testing
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentService) GetCommentByID(ctx context.Context, id uint) (*models.Comment, error) {
	args := m.Called(ctx, id)
	if comment, ok := args.Get(0).(*models.Comment); ok {
		return comment, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCommentService) GetCommentsByTask(ctx context.Context, taskID uint, page, pageSize int) ([]models.Comment, int64, error) {
	args := m.Called(ctx, taskID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentService) DeleteComment(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test CreateComment
func TestCreateComment(t *testing.T) {
	testCases := []struct {
		name           string
		inputComment   models.Comment
		mockSetup      func(*MockCommentService)
		expectedStatus int
	}{
		{
			name: "Successful Comment Creation",
			inputComment: models.Comment{
				TaskID:  1,
				Content: "Test comment",
				UserID:  1,
			},
			mockSetup: func(mcs *MockCommentService) {
				mcs.On("CreateComment", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid Input",
			inputComment: models.Comment{
				TaskID: 0, // Invalid TaskID
			},
			mockSetup: func(mcs *MockCommentService) {
				mcs.On("CreateComment", mock.Anything, mock.Anything).Return(errors.New("invalid input"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockCommentService)
			tc.mockSetup(mockService)

			// Create handler
			handler := NewCommentHandler(mockService)

			// Prepare request body
			jsonBody, _ := json.Marshal(tc.inputComment)
			req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Call handler
			handler.CreateComment(w, req)

			// Check response
			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

// Test GetCommentByID
func TestGetCommentByID(t *testing.T) {
	testCases := []struct {
		name           string
		commentID      string
		mockSetup      func(*MockCommentService)
		expectedStatus int
	}{
		{
			name:      "Successful Get Comment",
			commentID: "1",
			mockSetup: func(mcs *MockCommentService) {
				comment := &models.Comment{
					ID:      1,
					Content: "Test comment",
					TaskID:  1,
				}
				mcs.On("GetCommentByID", mock.Anything, uint(1)).Return(comment, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Invalid Comment ID",
			commentID: "0",
			mockSetup: func(mcs *MockCommentService) {
				mcs.On("GetCommentByID", mock.Anything, mock.Anything).Return(nil, errors.New("invalid ID"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Comment Not Found",
			commentID: "999",
			mockSetup: func(mcs *MockCommentService) {
				mcs.On("GetCommentByID", mock.Anything, uint(999)).Return(nil, errors.New("comment not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockCommentService)
			tc.mockSetup(mockService)

			// Create handler
			handler := NewCommentHandler(mockService)

			// Prepare request
			req := httptest.NewRequest(http.MethodGet, "/comments/"+tc.commentID, nil)
			req.SetPathValue("id", tc.commentID)
			w := httptest.NewRecorder()

			// Call handler
			handler.GetCommentByID(w, req)

			// Check response
			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

// Test GetCommentsByTask
// func TestGetCommentsByTask(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		taskID         string
// 		page           string
// 		pageSize       string
// 		mockSetup      func(*MockCommentService)
// 		expectedStatus int
// 	}{
// 		{
// 			name:     "Successful Get Comments",
// 			taskID:   "1",
// 			page:     "1",
// 			pageSize: "10",
// 			mockSetup: func(mcs *MockCommentService) {
// 				comments := []*models.Comment{
// 					{ID: 1, Content: "Comment 1", TaskID: 1},
// 					{ID: 2, Content: "Comment 2", TaskID: 1},
// 				}
// 				mcs.On("GetCommentsByTask", mock.Anything, uint(1), 1, 10).Return(comments, 2, nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:     "Default Pagination",
// 			taskID:   "1",
// 			page:     "",
// 			pageSize: "",
// 			mockSetup: func(mcs *MockCommentService) {
// 				comments := []*models.Comment{
// 					{ID: 1, Content: "Comment 1", TaskID: 1},
// 				}
// 				mcs.On("GetCommentsByTask", mock.Anything, uint(1), 1, 10).Return(comments, 1, nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:     "Service Error",
// 			taskID:   "1",
// 			page:     "1",
// 			pageSize: "10",
// 			mockSetup: func(mcs *MockCommentService) {
// 				mcs.On("GetCommentsByTask", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					Return(nil, 0, errors.New("service error"))
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Create mock service
// 			mockService := new(MockCommentService)
// 			tc.mockSetup(mockService)

// 			// Create handler
// 			handler := NewCommentHandler(mockService)

// 			// Prepare request
// 			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tc.taskID+"/comments", nil)
// 			req.SetPathValue("task_id", tc.taskID)
// 			if tc.page != "" {
// 				req.URL.Query().Set("page", tc.page)
// 			}
// 			if tc.pageSize != "" {
// 				req.URL.Query().Set("page_size", tc.pageSize)
// 			}
// 			w := httptest.NewRecorder()

// 			// Call handler
// 			handler.GetCommentsByTask(w, req)

// 			// Check response
// 			resp := w.Result()
// 			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

// 			// Verify mock expectations
// 			mockService.AssertExpectations(t)
// 		})
// 	}
// }

// Test DeleteComment
func TestDeleteComment(t *testing.T) {
	testCases := []struct {
		name           string
		commentID      string
		mockSetup      func(*MockCommentService)
		expectedStatus int
	}{
		{
			name:      "Successful Delete",
			commentID: "1",
			mockSetup: func(mcs *MockCommentService) {
				mcs.On("DeleteComment", mock.Anything, uint(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		// {
		// 	name:      "Invalid Comment ID",
		// 	commentID: "0",
		// 	mockSetup: func(mcs *MockCommentService) {
		// 		mcs.On("DeleteComment", mock.Anything, mock.Anything).Return(errors.New("invalid ID"))
		// 	},
		// 	expectedStatus: http.StatusBadRequest,
		// },
		{
			name:      "Delete Failed",
			commentID: "1",
			mockSetup: func(mcs *MockCommentService) {
				mcs.On("DeleteComment", mock.Anything, uint(1)).Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockCommentService)
			tc.mockSetup(mockService)

			// Create handler
			handler := NewCommentHandler(mockService)

			// Prepare request
			req := httptest.NewRequest(http.MethodDelete, "/comments/"+tc.commentID, nil)
			req.SetPathValue("id", tc.commentID)
			w := httptest.NewRecorder()

			// Call handler
			handler.DeleteComment(w, req)

			// Check response
			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
