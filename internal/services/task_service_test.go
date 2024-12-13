package services

import (
	"context"
	"example/project-management-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
    mock.Mock
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
    args := m.Called(ctx, task)
    return args.Error(0)
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error) {
    args := m.Called(ctx, projectID, page, pageSize)
    return args.Get(0).([]models.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
    args := m.Called(ctx, task)
    return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}


func TestCreateTask(t *testing.T) {
    t.Parallel()

    // Test cases
    testCases := []struct {
        name           string
        task           *models.Task
        mockRepoReturn error
        expectedError  bool
    }{
        // {
        //     name: "Successful Task Creation",
        //     task: &models.Task{
        //         Title:       "Test Task",
        //         Description: "Test Description",
        //         ProjectID:   1,
        //     },
        //     mockRepoReturn: nil,
        //     expectedError:  false,
        // },
        // {
        //     name: "task title is required",
        //     task: &models.Task{
        //         Title:       "",
        //         ProjectID:   1,
        //     },
        //     mockRepoReturn: errors.New("task title is required"),
        //     expectedError:  true,
        // },
        // {
        //     name: "No Project ID",
        //     task: &models.Task{
        //         Title: "Test Task",
        //     },
        //     mockRepoReturn: nil,
        //     expectedError:  true,
        // },
        {
            name: "Repository Error",
            task: &models.Task{
                Title:       "Test Task",
                ProjectID:   1,
            },
            mockRepoReturn: assert.AnError,
            expectedError:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mock repository
            mockRepo := new(MockTaskRepository)
            
            // Setup expectations
            mockRepo.On("CreateTask", mock.Anything, tc.task).Return(tc.mockRepoReturn)

            // Create service with mock repository
            service := NewTaskService(mockRepo)

            // Perform the test
            err := service.CreateTask(context.Background(), tc.task)

			// Assertions
            if tc.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            // Verify mock expectations
            if tc.name == "task title is required" || tc.name == "No Project ID" {
				mockRepo.AssertNumberOfCalls(t, "CreateTask", 0)
			} else {
				mockRepo.AssertExpectations(t)
			}
        })
    }
}

func TestGetTaskByID(t *testing.T) {
    t.Parallel()

    // Prepare test data
    expectedTask := &models.Task{
        BaseModel: models.BaseModel{ID: 1},
        Title:     "Test Task",
        ProjectID: 1,
    }

    testCases := []struct {
        name           string
        taskID         uint
        mockRepoReturn *models.Task
        mockRepoError  error
        expectedError  bool
    }{
        {
            name:           "Successful Task Retrieval",
            taskID:         1,
            mockRepoReturn: expectedTask,
            mockRepoError:  nil,
            expectedError:  false,
        },
        {
            name:           "Task Not Found",
            taskID:         999,
            mockRepoReturn: &models.Task{},
            mockRepoError:  assert.AnError,
            expectedError:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mock repository
            mockRepo := new(MockTaskRepository)
            
            // Setup expectations
            mockRepo.On("GetTaskByID", mock.Anything, tc.taskID).
                Return(tc.mockRepoReturn, tc.mockRepoError)

            // Create service with mock repository
            service := NewTaskService(mockRepo)

            // Perform the test
            task, err := service.GetTaskByID(context.Background(), tc.taskID)

            // Assertions
            if tc.expectedError {
                assert.Error(t, err)
                assert.Nil(t, task)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, expectedTask, task)
            }

            // Verify mock expectations
            mockRepo.AssertExpectations(t)
        })
    }
}

func TestUpdateTask(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name           string
        task           *models.Task
        mockRepoReturn error
        expectedError  bool
    }{
        {
            name: "Successful Task Update",
            task: &models.Task{
                BaseModel: models.BaseModel{ID: 1},
                Title:     "Updated Task",
                ProjectID: 1,
            },
            mockRepoReturn: nil,
            expectedError:  false,
        },
        {
            name: "Empty Task Title",
            task: &models.Task{
                BaseModel: models.BaseModel{ID: 1},
                Title:     "",
                ProjectID: 1,
            },
            mockRepoReturn: nil,
            expectedError:  true,
        },
        {
            name: "Repository Error",
            task: &models.Task{
                BaseModel: models.BaseModel{ID: 1},
                Title:     "Updated Task",
                ProjectID: 1,
            },
            mockRepoReturn: assert.AnError,
            expectedError:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mock repository
            mockRepo := new(MockTaskRepository)
            
            // Setup expectations
            mockRepo.On("UpdateTask", mock.Anything, tc.task).Return(tc.mockRepoReturn)

            // Create service with mock repository
            service := NewTaskService(mockRepo)

            // Perform the test
            err := service.UpdateTask(context.Background(), tc.task)

            // Assertions
            if tc.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            // Verify mock expectations
			if tc.name == "Empty Task Title" {
				mockRepo.AssertNotCalled(t, "UpdateTask")
			} else {
				mockRepo.AssertExpectations(t)
			}
        })
    }
}

func TestGetTasksByProject(t *testing.T) {
    t.Parallel()

    // Prepare test data
    projectID := uint(1)
    expectedTasks := []models.Task{
        {
            BaseModel: models.BaseModel{ID: 1},
            Title:     "Task 1",
            ProjectID: projectID,
        },
        {
            BaseModel: models.BaseModel{ID: 2},
            Title:     "Task 2",
            ProjectID: projectID,
        },
    }
    expectedTotal := int64(2)

    testCases := []struct {
        name           string
        projectID      uint
        page           int
        pageSize       int
        mockTasksReturn []models.Task
        mockTotalReturn int64
        mockRepoError   error
        expectedError   bool
    }{
        {
            name:           "Successful Tasks Retrieval",
            projectID:      projectID,
            page:           1,
            pageSize:       10,
            mockTasksReturn: expectedTasks,
            mockTotalReturn: expectedTotal,
            mockRepoError:   nil,
            expectedError:   false,
        },
        {
            name:           "Repository Error",
            projectID:      projectID,
            page:           1,
            pageSize:       10,
            mockTasksReturn: nil,
            mockTotalReturn: 0,
            mockRepoError:   assert.AnError,
            expectedError:   true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mock repository
            mockRepo := new(MockTaskRepository)
            
            // Setup expectations
            mockRepo.On("GetTaskByProject", 
                mock.Anything, 
                tc.projectID, 
                tc.page, 
                tc.pageSize,
            ).Return(tc.mockTasksReturn, tc.mockTotalReturn, tc.mockRepoError)

            // Create service with mock repository
            service := NewTaskService(mockRepo)

            // Perform the test
            tasks, total, err := service.GetTasksByProject(
                context.Background(), 
                tc.projectID, 
                tc.page, 
                tc.pageSize,
            )

            // Assertions
            if tc.expectedError {
                assert.Error(t, err)
                assert.Nil(t, tasks)
                assert.Zero(t, total)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tc.mockTasksReturn, tasks)
                assert.Equal(t, tc.mockTotalReturn, total)
            }

            // Verify mock expectations
            mockRepo.AssertExpectations(t)
        })
    }
}

func TestDeleteTask(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name           string
        taskID         uint
        mockRepoReturn error
        expectedError  bool
    }{
        {
            name:           "Successful Task Deletion",
            taskID:         1,
            mockRepoReturn: nil,
            expectedError:  false,
        },
        {
            name:           "Repository Error",
            taskID:         1,
            mockRepoReturn: assert.AnError,
            expectedError:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mock repository
            mockRepo := new(MockTaskRepository)
            
            // Setup expectations
            mockRepo.On("DeleteTask", mock.Anything, tc.taskID).Return(tc.mockRepoReturn)

            // Create service with mock repository
            service := NewTaskService(mockRepo)

            // Perform the test
            err := service.DeleteTask(context.Background(), tc.taskID)

            // Assertions
            if tc.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            // Verify mock expectations
            mockRepo.AssertExpectations(t)
        })
    }
}
