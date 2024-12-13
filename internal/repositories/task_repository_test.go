package repositories

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"example/project-management-system/internal/models"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		DriverName:                "mysql",
		SkipInitializeWithVersion: true,
	})
	
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

func TestCreateTask(t *testing.T) {
	// Prepare test cases
	testCases := []struct {
		name           string
		task           *models.Task
		mockExpectFunc func(mock sqlmock.Sqlmock)
		expectedError  bool
	}{
		{
			name: "Successful Task Creation",
			task: &models.Task{
				Title:       "Test Task",
				Description: "Test Description",
				ProjectID:   1,
				AssignedTo:  2,
			},
			mockExpectFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tasks`")).
					WithArgs(
						sqlmock.AnyArg(),  // CreatedAt
						sqlmock.AnyArg(),  // UpdatedAt
						nil,  // DeletedAt
						"Test Task",
						"Test Description",
						uint(1),  // ProjectID
						uint(2),  // AssignedTo
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Database Error",
			task: &models.Task{
				Title:       "Test Task",
				ProjectID:   1,
			},
			mockExpectFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tasks`")).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			db, mock := setupMockDB(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			// Prepare mock expectations
			tc.mockExpectFunc(mock)

			// Create repository
			repo := NewTaskRepository(db)

			// Execute create task
			err := repo.CreateTask(context.Background(), tc.task)

			// Assertions
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tc.task.ID)
			}

			// Verify mock expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// func TestGetTaskByID(t *testing.T) {
// 	// Prepare test data
// 	now := time.Now()
// 	testTask := models.Task{
// 		BaseModel: models.BaseModel{
// 			ID:        1,
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		},
// 		Title:       "Test Task",
// 		Description: "Test Description",
// 		ProjectID:   1,
// 		AssignedTo:  2,
// 		Project: models.Project{
// 			BaseModel: models.BaseModel{ID: 1},
// 			Name:      "Test Project",
// 		},
// 		Assignee: models.User{
// 			BaseModel: models.BaseModel{ID: 2},
// 			Username:  "testuser",
// 		},
// 	}

// 	testCases := []struct {
// 		name           string
// 		taskID         uint
// 		mockExpectFunc func(mock sqlmock.Sqlmock)
// 		expectedError  bool
// 	}{
// 		{
// 			name:   "Successful Task Retrieval",
// 			taskID: 1,
// 			mockExpectFunc: func(mock sqlmock.Sqlmock) {
// 				// Expect task query
// 				rows := sqlmock.NewRows([]string{
// 					"id", "created_at", "updated_at", "deleted_at", 
// 					"title", "description", "project_id", "assigned_to",
// 				}).AddRow(
// 					testTask.ID, testTask.CreatedAt, testTask.UpdatedAt, nil,
// 					testTask.Title, testTask.Description, testTask.ProjectID, testTask.AssignedTo,
// 				)
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks`")).
// 					WithArgs(testTask.ID, 1).
// 					WillReturnRows(rows)

// 				// Expect project preload
// 				projectRows := sqlmock.NewRows([]string{
// 					"id", "name",
// 				}).AddRow(
// 					testTask.Project.ID, testTask.Project.Name,
// 				)
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects`")).
// 					WithArgs(testTask.ProjectID, 1).
// 					WillReturnRows(projectRows)

// 				// Expect assignee preload
// 				assigneeRows := sqlmock.NewRows([]string{
// 					"id", "username",
// 				}).AddRow(
// 					testTask.Assignee.ID, testTask.Assignee.Username,
// 				)
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
// 					WithArgs(testTask.AssignedTo, 1).
// 					WillReturnRows(assigneeRows)
// 			},
// 			expectedError: false,
// 		},
// 		// {
// 		// 	name:   "Task Not Found",
// 		// 	taskID: 999,
// 		// 	mockExpectFunc: func(mock sqlmock.Sqlmock) {
// 		// 		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks`")).
// 		// 			WithArgs(uint(999)).
// 		// 			WillReturnError(gorm.ErrRecordNotFound)
// 		// 	},
// 		// 	expectedError: true,
// 		// },
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Setup mock database
// 			db, mock := setupMockDB(t)
// 			defer func() {
// 				sqlDB, _ := db.DB()
// 				sqlDB.Close()
// 			}()

// 			// Prepare mock expectations
// 			tc.mockExpectFunc(mock)

// 			// Create repository
// 			repo := NewTaskRepository(db)

// 			// Execute get task by ID
// 			task, err := repo.GetTaskByID(context.Background(), tc.taskID)

// 			// Assertions
// 			if tc.expectedError {
// 				assert.Error(t, err)
// 				assert.Nil(t, task)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, task)
// 				assert.Equal(t, testTask.Title, task.Title)
// 			}

// 			// Verify mock expectations
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

func TestUpdateTask(t *testing.T) {
	testCases := []struct {
		name           string
		task           *models.Task
		mockExpectFunc func(mock sqlmock.Sqlmock)
		expectedError  bool
	}{
		{
			name: "Successful Task Update",
			task: &models.Task{
				BaseModel: models.BaseModel{ID: 1},
				Title:     "Updated Task",
				Description: "Updated Description",
				ProjectID:   1,
				AssignedTo:  2,
			},
			mockExpectFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks`")).
					WithArgs(
						sqlmock.AnyArg(), // UpdatedAt,
						sqlmock.AnyArg(),
						nil,
						"Updated Task",
						"Updated Description",
						uint(1),  // ProjectID
						uint(2),  // AssignedTo
						1,        // ID
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Update Error",
			task: &models.Task{
				BaseModel: models.BaseModel{ID: 1},
				Title:     "Updated Task",
			},
			mockExpectFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks`")).
					WillReturnError(errors.New("database update error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			db, mock := setupMockDB(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			// Prepare mock expectations
			tc.mockExpectFunc(mock)

			// Create repository
			repo := NewTaskRepository(db)

			// Execute update task
			err := repo.UpdateTask(context.Background(), tc.task)

			// Assertions
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteTask(t *testing.T) {
	testCases := []struct {
		name           string
		taskID         uint
		mockExpectFunc func(mock sqlmock.Sqlmock)
		expectedError  bool
	}{
		{
			name:   "Successful Task Deletion",
			taskID: 1,
			mockExpectFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks`")).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:   "Delete Error",
			taskID: 1,
			mockExpectFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks`")).
					WillReturnError(errors.New("database delete error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			db, mock := setupMockDB(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			// Prepare mock expectations
			tc.mockExpectFunc(mock)

			// Create repository
			repo := NewTaskRepository(db)

			// Execute delete task
			err := repo.DeleteTask(context.Background(), tc.taskID)

			// Assertions
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// func TestGetTaskByProject(t *testing.T) {
// 	// Prepare test data
// 	projectID := uint(1)
// 	testTasks := []models.Task{
// 		{
// 			BaseModel: models.BaseModel{ID: 1},
// 			Title:     "Task 1",
// 			ProjectID: projectID,
// 		},
// 		{
// 			BaseModel: models.BaseModel{ID: 2},
// 			Title:     "Task 2",
// 			ProjectID: projectID,
// 		},
// 	}

// 	testCases := []struct {
// 		name           string
// 		projectID      uint
// 		page           int
// 		pageSize       int
// 		mockExpectFunc func(mock sqlmock.Sqlmock)
// 		expectedError  bool
// 	}{
// 		{
// 			name:      "Successful Tasks Retrieval",
// 			projectID: projectID,
// 			page:      1,
// 			pageSize:  10,
// 			mockExpectFunc: func(mock sqlmock.Sqlmock) {
// 				// Mock count query
// 				countRows := sqlmock.NewRows([]string{"count"}).
// 					AddRow(int64(len(testTasks)))
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `tasks`")).
// 					WithArgs(projectID).
// 					WillReturnRows(countRows)

// 				// Mock tasks query
// 				tasksRows := sqlmock.NewRows([]string{
// 					"id", "title", "project_id",
// 				})
// 				for _, task := range testTasks {
// 					tasksRows.AddRow(
// 						task.ID, 
// 						task.Title, 
// 						task.ProjectID,
// 					)
// 				}
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks`")).
// 					WithArgs(projectID).
// 					WillReturnRows(tasksRows)
// 			},
// 			expectedError: false,
// 		},
// 		{
// 			name:      "Database Error",
// 			projectID: projectID,
// 			page:      1,
// 			pageSize:  10,
// 			mockExpectFunc: func(mock sqlmock.Sqlmock) {
// 				// Mock count query error
// 				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `tasks`")).
// 					WithArgs(projectID).
// 					WillReturnError(errors.New("database error"))
// 			},
// 			expectedError: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Setup mock database
// 			db, mock := setupMockDB(t)
// 			defer func() {
// 				sqlDB, _ := db.DB()
// 				sqlDB.Close()
// 			}()

// 			// Prepare mock expectations
// 			tc.mockExpectFunc(mock)

// 			// Create repository
// 			repo := NewTaskRepository(db)

// 			// Execute get tasks by project
// 			tasks, total, err := repo.GetTaskByProject(
// 				context.Background(), 
// 				tc.projectID, 
// 				tc.page, 
// 				tc.pageSize,
// 			)

// 			// Assertions
// 			if tc.expectedError {
// 				assert.Error(t, err)
// 				assert.Nil(t, tasks)
// 				assert.Zero(t, total)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Len(t, tasks, len(testTasks))
// 				assert.Equal(t, int64(len(testTasks)), total)
// 			}
// 			// Verify mock expectations
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }
