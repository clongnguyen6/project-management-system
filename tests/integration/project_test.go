package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"example/project-management-system/internal/config"
	"example/project-management-system/internal/database"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/server"
	"example/project-management-system/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB   *gorm.DB
	testCtx  context.Context
	testLogger logger.Logger
	testServer *httptest.Server

)

func insertTestProduct(name string, description string) (int, error) {
	var id int
	query := "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id"
	err := testDB.Raw(query, name, description).Preload("tasks").Preload("teams").Scan(&id)
	return id, err.Error
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Container request for PostgreSQL
	req := testcontainers.ContainerRequest{
		Image: "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	// Start PostgreSQL container
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		if err := postgresC.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}()

	// Get container connection details
	host, err := postgresC.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}

	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get mapped port: %v", err)
	}

	// Create test config
	testConfig := &config.Config{
		Addr:  "localhost:8080",
		Db: config.DbConfig{
			Host:     host,
			Port:     port.Port(),
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
		AUTH0_DOMAIN: "https://dev-n5mocwlrk8i63cjm.us.auth0.com/",
		AUTH0_AUDIENCE: "https://project-management-api",
		ENVIRONMENT:	"test",
	}

	dbInstance := database.NewPostgresConnection(testConfig, testLogger)
	testDB = dbInstance.DB

	// Run migrations
	err = testDB.AutoMigrate(&models.Project{}, &models.Task{}, &models.Team{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	httpServer, err := server.NewHTTPServer(testDB, testConfig)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	testServer = httptest.NewServer(httpServer.Handler)

	// Run tests
	code := m.Run()

	// Cleanup
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	testServer.Close()

	os.Exit(code)
}

func TestCreateProject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		project := map[string]interface{}{
			"name":        "Test Project",
			"description": "This is a test project",
		}
		body, err := json.Marshal(project)
		if err != nil {
			t.Fatalf("Failed to marshal project: %v", err)
		}
	
		req, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/v1/projects", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, "Test Project", response["name"])

	})
}

func TestGetAllProjects(t *testing.T) {
	t.Run("Get All Projects", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/api/v1/projects", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, []interface{}{}, response["projects"])

	})
}


func TestGetProjectByID(t *testing.T) {
	id, _ := insertTestProduct("Test Project 1", "Decription Test Project 1")
	t.Run("Get project by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+fmt.Sprintf("/api/v1/projects/%d", id), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, "Test Project 1", response["name"])

	})
}

func TestUpdateProject(t *testing.T) {
	id, _ := insertTestProduct("Test Project 1", "Decription Test Project 1")
	t.Run("Update Project", func(t *testing.T) {

		project := map[string]interface{}{
			"name":        "Updated Project",
		}
		body, err := json.Marshal(project)
		if err != nil {
			t.Fatalf("Failed to marshal project: %v", err)
		}

		req, err := http.NewRequest(http.MethodPut, testServer.URL+fmt.Sprintf("/api/v1/projects/%d", id), bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, "Updated Project", response["name"])

	})
}

func TestDeleteProject(t *testing.T) {
	id, _ := insertTestProduct("Test Project 1", "Decription Test Project 1")
	t.Run("Delete Project", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, testServer.URL+fmt.Sprintf("/api/v1/projects/%d", id), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, nil, response["name"])

	})
}
