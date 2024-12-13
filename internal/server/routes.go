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
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userHandler.CreateUser),
	)
	router.HandleFunc("GET /api/v1/users",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userHandler.GetAllUsers),
	)
	router.HandleFunc("GET /api/v1/users/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userHandler.GetUserByID),
	)
	router.HandleFunc("DELETE /api/v1/users/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userHandler.DeleteUser),
	)
	router.HandleFunc("POST /api/v1/projects/{projectId}/users/{userId}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userProjectHandler.AddUserToProject),
	)

	router.HandleFunc("POST /api/v1/projects",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, projectHandler.CreateProject),
	)
	router.HandleFunc("GET /api/v1/projects",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, projectHandler.GetAllProjects),
	)
	router.HandleFunc("GET /api/v1/projects/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, projectHandler.GetProjectByID),
	)
	router.HandleFunc("PUT /api/v1/projects/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, projectHandler.UpdateProject),
	)
	router.HandleFunc("DELETE /api/v1/projects/{id}",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, projectHandler.DeleteProject),
	)
	router.HandleFunc("GET /api/v1/projects/{projectID}/tasks",
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, projectHandler.GetTaskByProjectID),
	)


	router.HandleFunc("POST /api/v1/tasks", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, taskHandler.CreateTask),
	)
	router.HandleFunc("GET /api/v1/tasks/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, taskHandler.GetTaskByID),
	)
	router.HandleFunc("GET /api/v1/tasks", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, taskHandler.GetTasksByProject),
	)
	router.HandleFunc("PUT /api/v1/tasks/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, taskHandler.UpdateTask),
	)
	router.HandleFunc("DELETE /api/v1/tasks/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, taskHandler.DeleteTask),
	)


	router.HandleFunc("POST /api/v1/teams", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, teamHandler.CreateTeam),
	)
	router.HandleFunc("GET /api/v1/teams/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, teamHandler.GetTeamByID),
	)
	router.HandleFunc("GET /api/v1/teams", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, teamHandler.GetPaginatedTeams),
	)
	router.HandleFunc("PUT /api/v1/teams/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, teamHandler.UpdateTeam),
	)
	router.HandleFunc("DELETE /api/v1/teams/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, teamHandler.DeleteTeam),
	)


	router.HandleFunc("POST /api/v1/comments", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, commentHandler.CreateComment),
	)
	router.HandleFunc("GET /api/v1/comments/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, commentHandler.GetCommentByID),
	)
	router.HandleFunc("GET /api/v1/comments", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, commentHandler.GetCommentsByTask),
	)
	router.HandleFunc("DELETE /api/v1/comments/{id}", 
		middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, commentHandler.DeleteComment),
	)


	// router.HandleFunc("POST /api/v1/users-projects/{projectId}/users/{userId}", 
	// 	middleware.ValidateJWT(cfg.AUTH0_AUDIENCE, cfg.AUTH0_DOMAIN, userProjectHandler.AddUserToProject),
	// )
	
	// return middleware.HandleCacheControl(router)
	return router
}

