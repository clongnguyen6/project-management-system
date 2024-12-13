package repositories

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"gorm.io/gorm"

// 	"example/project-management-system/internal/models"
// )

// // MockTaskRepository mô phỏng TaskRepository
// type MockTaskRepository struct {
// 	mock.Mock
// }

// func (m *MockTaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
// 	args := m.Called(ctx, task)
// 	return args.Error(0)
// }

// func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
// 	args := m.Called(ctx, id)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.Task), args.Error(1)
// }

// func (m *MockTaskRepository) GetTaskByProject(ctx context.Context, projectID uint, page, pageSize int) ([]models.Task, int64, error) {
// 	args := m.Called(ctx, projectID, page, pageSize)
// 	return args.Get(0).([]models.Task), args.Get(1).(int64), args.Error(2)
// }

// func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
// 	args := m.Called(ctx, task)
// 	return args.Error(0)
// }

// func (m *MockTaskRepository) DeleteTask(ctx context.Context, id uint) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// // Integrated Repository Tests (Requires actual database setup)
// func TestTaskRepositoryIntegration(t *testing.T) {
// 	// Bạn cần setup database test ở đây
// 	// Sử dụng test database hoặc in-memory database như sqlite
// 	// db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
// 	// assert.NoError(t, err)
// 	// defer db.Close()

// 	// Migrate models
// 	// db.AutoMigrate(&models.Task{}, &models.Project{}, &models.User{})

// 	// repo := repositories.NewTaskRepository(db)

// 	t.Run("Create Task", func(t *testing.T) {
// 		// Setup mock
// 		mockRepo := new(MockTaskRepository)

// 		// Prepare test task
// 		testTask := &models.Task{
// 			Title:       "Test Task",
// 			Description: "Test Description",
// 			ProjectID:   1,
// 			AssignedTo:  1,
// 		}

// 		// Expect create method to be called
// 		mockRepo.On("CreateTask", mock.Anything, testTask).Return(nil)

// 		// Execute
// 		err := mockRepo.CreateTask(context.Background(), testTask)

// 		// Assert
// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Get Task By ID", func(t *testing.T) {
// 		mockRepo := new(MockTaskRepository)

// 		// Prepare mock task
// 		mockTask := &models.Task{
// 			BaseModel: models.BaseModel{
// 				ID:        1,
// 				CreatedAt: time.Now(),
// 				UpdatedAt: time.Now(),
// 			},
// 			Title:       "Existing Task",
// 			Description: "Task Description",
// 			ProjectID:   1,
// 			AssignedTo:  1,
// 		}

// 		// Setup expectation
// 		mockRepo.On("GetTaskByID", mock.Anything, uint(1)).Return(mockTask, nil)

// 		// Execute
// 		task, err := mockRepo.GetTaskByID(context.Background(), 1)

// 		// Assert
// 		assert.NoError(t, err)
// 		assert.NotNil(t, task)
// 		assert.Equal(t, mockTask.Title, task.Title)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Get Tasks By Project", func(t *testing.T) {
// 		mockRepo := new(MockTaskRepository)

// 		// Prepare mock tasks
// 		mockTasks := []models.Task{
// 			{
// 				BaseModel: models.BaseModel{ID: 1},
// 				Title:     "Task 1",
// 				ProjectID: 1,
// 			},
// 			{
// 				BaseModel: models.BaseModel{ID: 2},
// 				Title:     "Task 2",
// 				ProjectID: 1,
// 			},
// 		}

// 		// Setup expectation
// 		mockRepo.On("GetTaskByProject", mock.Anything, uint(1), 1, 10).
// 			Return(mockTasks, int64(2), nil)

// 		// Execute
// 		tasks, total, err := mockRepo.GetTaskByProject(context.Background(), 1, 1, 10)

// 		// Assert
// 		assert.NoError(t, err)
// 		assert.Len(t, tasks, 2)
// 		assert.Equal(t, int64(2), total)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Update Task", func(t *testing.T) {
// 		mockRepo := new(MockTaskRepository)

// 		// Prepare task to update
// 		updateTask := &models.Task{
// 			BaseModel: models.BaseModel{ID: 1},
// 			Title:     "Updated Task Title",
// 		}

// 		// Setup expectation
// 		mockRepo.On("UpdateTask", mock.Anything, updateTask).Return(nil)

// 		// Execute
// 		err := mockRepo.UpdateTask(context.Background(), updateTask)

// 		// Assert
// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Delete Task", func(t *testing.T) {
// 		mockRepo := new(MockTaskRepository)

// 		// Setup expectation
// 		mockRepo.On("DeleteTask", mock.Anything, uint(1)).Return(nil)

// 		// Execute
// 		err := mockRepo.DeleteTask(context.Background(), 1)

// 		// Assert
// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	// Error Scenarios
// 	t.Run("Get Non-Existent Task", func(t *testing.T) {
// 		mockRepo := new(MockTaskRepository)

// 		// Setup expectation for non-existent task
// 		mockRepo.On("GetTaskByID", mock.Anything, uint(999)).
// 			Return(nil, gorm.ErrRecordNotFound)

// 		// Execute
// 		task, err := mockRepo.GetTaskByID(context.Background(), 999)

// 		// Assert
// 		assert.Error(t, err)
// 		assert.Nil(t, task)
// 		assert.Equal(t, gorm.ErrRecordNotFound, err)
// 		mockRepo.AssertExpectations(t)
// 	})
// }

// // Nếu bạn muốn test với database thực tế
// // Cần setup database test environment
// func setupTestDatabase() *gorm.DB {
// 	// Implement database connection for testing
// 	// Sử dụng in-memory database hoặc test database
// 	return nil
// }