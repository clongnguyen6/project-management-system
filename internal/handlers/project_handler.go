package handlers

import (
	"encoding/json"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/services"
	"example/project-management-system/internal/utils/response"
	"fmt"
	"net/http"
	"strconv"
)

type ProjectHandler interface {
	CreateProject(w http.ResponseWriter, r *http.Request)
	GetProjectByID(w http.ResponseWriter, r *http.Request)
	GetAllProjects(w http.ResponseWriter, r *http.Request)
	UpdateProject(w http.ResponseWriter, r *http.Request)
	DeleteProject(w http.ResponseWriter, r *http.Request)
	GetTaskByProjectID(w http.ResponseWriter, r *http.Request)
}

type ProjectHandlerImplementation struct {
	service services.ProjectService
}

func NewProjectHandler(service services.ProjectService) *ProjectHandlerImplementation {
	return &ProjectHandlerImplementation{service: service}
}


// CreateProject godoc
//	@Summary		Create a new project
//	@Description	Create a new project with provided details
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			project	body		models.Project		true	"Project Creation Request"
//	@Success		201		{object}	map[string]int		"Project created successfully"
//	@Failure		400		{object}	response.Response	"Invalid input"
//	@Failure		500		{object}	response.Response	"Server error"
//	@Router			/projects [post]
func (h *ProjectHandlerImplementation) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteJson(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("%s", "Method not allowed")))
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.CreateProject(r.Context(), &project); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusCreated, project)
}


// GetUsers godoc
//	@Summary		Get all projects
//	@Description	Retrieve paginated list of projects
//	@Tags			Projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int						false	"Page number"					default(1)
//	@Param			pageSize	query		int						false	"Number of projects per page"	default(10)
//	@Success		200			{object}	map[string]interface{}	"Successful response"
//	@Failure		400			{object}	response.Response		"Bad request"
//	@Router			/projects [get]
func (h *ProjectHandlerImplementation) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, _ := strconv.Atoi(pageStr)
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize == 0 {
		pageSize = 10
	}

	projects, total, err := h.service.GetPaginatedProjects(r.Context(), page, pageSize)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"projects": projects,
		"total":    total,
		"page":     page,
	})
}


// GetUserByID godoc
//	@Summary		Get project by ID
//	@Description	Retrieve a specific project by their unique identifier
//	@Tags			Projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Project ID"
//	@Success		200	{object}	models.Project		"Successful response"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		404	{object}	response.Response	"User not found"
//	@Router			/projects/{id} [get]
func (h *ProjectHandlerImplementation) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteJson(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("%s", "Method not allowed")))
		return
	}

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	project, err := h.service.GetProjectByID(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, project)
}


// UpdateProject godoc
//	@Summary		Update an existing project
//	@Description	Update a project's details by its ID
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Project ID"
//	@Param			project	body		models.Project		true	"Updated Project Details"
//	@Success		200		{object}	models.Project		"Successful response"
//	@Failure		400		{object}	response.Response	"Bad request"
//	@Failure		404		{object}	response.Response	"User not found"
//	@Router			/projects/{id} [put]
func (h *ProjectHandlerImplementation) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	project.ID = uint(id)
	if err := h.service.UpdateProject(r.Context(), &project); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, project)
}


// DeleteUser godoc
//	@Summary		Delete a project
//	@Description	Remove a project from the system by their ID
//	@Tags			Projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Project ID to delete"
//	@Success		200	{object}	map[string]string	"Project deleted successfully"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		500	{object}	response.Response	"Server error"
//	@Router			/projects/{id} [delete]
func (h *ProjectHandlerImplementation) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.DeleteProject(r.Context(), uint(id)); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "project deleted successfully"})
}


// DeleteUser godoc
//	@Summary		Get task by project ID
//	@Description	Retrieve a list task by their ID
//	@Tags			Projects
//	@Produce		json
//	@Security		BearerAuth
//	@Param			projectID	path		int					true	"Project ID"
//	@Success		200			{object}	map[string]string	"Successful response"
//	@Failure		400			{object}	response.Response	"Bad request"
//	@Failure		500			{object}	response.Response	"Server error"
//	@Router			/projects/{projectID}/tasks [get]
func (h *ProjectHandlerImplementation) GetTaskByProjectID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("projectID"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	tasks, err := h.service.GetTaskByProjectID(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, tasks)
}
