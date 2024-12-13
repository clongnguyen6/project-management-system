package server

import (
	"example/project-management-system/internal/config"
	"example/project-management-system/internal/handlers"
	"example/project-management-system/internal/repositories"
	"example/project-management-system/internal/services"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func NewHTTPServer(db *gorm.DB, cfg *config.Config) (*http.Server, error) {
	router := http.NewServeMux()

	// Set up the api repositories
	userRepository := repositories.NewUserRepository(db)
	projectRepository := repositories.NewProjectRepository(db)
	taskRepository := repositories.NewTaskRepository(db)
	teamRepository := repositories.NewTeamRepository(db)
	commentRepository := repositories.NewCommentRepository(db)
	userProjectRepository := repositories.NewUserProjectRepository(db)

	// Set up the api services
	userService := services.NewUserService(userRepository)
	projectService := services.NewProjectService(projectRepository)
	taskService := services.NewTaskService(taskRepository)
	teamService := services.NewTeamService(teamRepository)
	commentService := services.NewCommentService(commentRepository)
	userProjectService := services.NewUserProjectService(userProjectRepository)

	// Set up the api handlers
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	teamHandler := handlers.NewTeamHandler(teamService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userProjectHandler := handlers.NewUserProjectHandler(userProjectService)

	// Set up API routes
	handler := RegisterRoutes(
		cfg,
		router,
		userHandler,
		projectHandler,
		taskHandler,
		teamHandler,
		commentHandler,
		userProjectHandler,
	)

	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Declare Server config
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, nil
}
