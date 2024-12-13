package server

import (
	"example/project-management-system/internal/config"
	"example/project-management-system/internal/handlers"
	"example/project-management-system/pkg/middleware"
	"net/http"
)


func RegisterRoutes(
	cfg *config.Config,
	router *http.ServeMux,
	userHandler handlers.UserHandler,
	projectHandler handlers.ProjectHandler,
	taskHandler handlers.TaskHandler,
	teamHandler handlers.TeamHandler,
	commentHandler handlers.CommentHandler,
	userProjectHandler handlers.UserProjectHandler,
) http.Handler {

	router.HandleFunc("POST /api/v1/users",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, userHandler.CreateUser),
	)
	router.HandleFunc("GET /api/v1/users",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, userHandler.GetAllUsers),
	)
	router.HandleFunc("GET /api/v1/users/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, userHandler.GetUserByID),
	)
	router.HandleFunc("DELETE /api/v1/users/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, userHandler.DeleteUser),
	)
	router.HandleFunc("POST /api/v1/projects/{projectId}/users/{userId}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, userProjectHandler.AddUserToProject),
	)

	router.HandleFunc("POST /api/v1/projects",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, projectHandler.CreateProject),
	)
	// router.HandleFunc("POST /api/v1/projects",
	// 	projectHandler.CreateProject,
	// )
	router.HandleFunc("GET /api/v1/projects",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, projectHandler.GetAllProjects),
	)
	router.HandleFunc("GET /api/v1/projects/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, projectHandler.GetProjectByID),
	)
	router.HandleFunc("PUT /api/v1/projects/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, projectHandler.UpdateProject),
	)
	router.HandleFunc("DELETE /api/v1/projects/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, projectHandler.DeleteProject),
	)
	router.HandleFunc("GET /api/v1/projects/{projectID}/tasks",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, projectHandler.GetTaskByProjectID),
	)


	router.HandleFunc("POST /api/v1/tasks", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, taskHandler.CreateTask),
	)
	router.HandleFunc("GET /api/v1/tasks/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, taskHandler.GetTaskByID),
	)
	router.HandleFunc("GET /api/v1/tasks", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, taskHandler.GetTasksByProject),
	)
	router.HandleFunc("PUT /api/v1/tasks/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, taskHandler.UpdateTask),
	)
	router.HandleFunc("DELETE /api/v1/tasks/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, taskHandler.DeleteTask),
	)


	router.HandleFunc("POST /api/v1/teams", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, teamHandler.CreateTeam),
	)
	router.HandleFunc("GET /api/v1/teams/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, teamHandler.GetTeamByID),
	)
	router.HandleFunc("GET /api/v1/teams", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, teamHandler.GetPaginatedTeams),
	)
	router.HandleFunc("PUT /api/v1/teams/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, teamHandler.UpdateTeam),
	)
	router.HandleFunc("DELETE /api/v1/teams/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, teamHandler.DeleteTeam),
	)


	router.HandleFunc("POST /api/v1/comments", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, commentHandler.CreateComment),
	)
	router.HandleFunc("GET /api/v1/comments/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, commentHandler.GetCommentByID),
	)
	router.HandleFunc("GET /api/v1/comments", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, commentHandler.GetCommentsByTask),
	)
	router.HandleFunc("DELETE /api/v1/comments/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, cfg.ENVIRONMENT, commentHandler.DeleteComment),
	)


	// router.HandleFunc("POST /api/v1/users-projects/{projectId}/users/{userId}", 
	// 	middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userProjectHandler.AddUserToProject),
	// )
	
	// return middleware.HandleCacheControl(router)
	return router
}

