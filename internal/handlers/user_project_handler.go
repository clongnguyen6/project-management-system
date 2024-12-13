package handlers

import (
	"example/project-management-system/internal/services"
	"example/project-management-system/internal/utils/response"
	"fmt"
	"net/http"
	"strconv"
)


type UserProjectHandler interface {
	AddUserToProject(w http.ResponseWriter, r *http.Request)
}

type UserProjectImplementation struct {
	userProjectService services.UserProjectService
}

// NewProjectHandler creates a new instance of the UserProjectImplementation.
func NewUserProjectHandler(userProjectService services.UserProjectService) *UserProjectImplementation {
	return &UserProjectImplementation{
		userProjectService: userProjectService,
	}
}

// AddUserToProject adds a user to a project.
//	@Summary		Add a user to a project
//	@Description	Add a user to a specified project by their IDs
//	@Tags			user_project
//	@Security		BearerAuth
//	@Param			userId		path	int	true	"User ID"
//	@Param			projectId	path	int	true	"Project ID"
//	@Success		204
//	@Router			/users-projects/{projectId}/users/{userId} [post]
func (handler *UserProjectImplementation) AddUserToProject(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(r.PathValue("userId"), 10, 64)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("%s", "Invalid user ID")))
		return
	}

	projectId, err := strconv.ParseUint(r.PathValue("projectId"), 10, 64)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("%s", "Invalid project ID")))
		return
	}


	if err := handler.userProjectService.AddUserToProject(uint(userId), uint(projectId)); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("%s", err.Error())))
		return
	}

	response.WriteJson(w, http.StatusNoContent, response.StatusOK)
}
