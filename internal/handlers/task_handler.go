package handlers

import (
	"encoding/json"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/services"
	"example/project-management-system/internal/utils/response"
	"net/http"

	"strconv"
)

type TaskHandler interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	GetTasksByProject(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

type TaskHandlerImplementation struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandlerImplementation {
	return &TaskHandlerImplementation{service: service}
}

// CreateProject godoc
//	@Summary		Create a new task
//	@Description	Create a new task with provided details
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			task	body		models.Task			true	"Task Creation Request"
//	@Success		201		{object}	map[string]int		"Task created successfully"
//	@Failure		400		{object}	response.Response	"Invalid input"
//	@Failure		500		{object}	response.Response	"Server error"
//	@Router			/tasks [post]
func (h *TaskHandlerImplementation) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}
	if err := h.service.CreateTask(r.Context(), &task); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusCreated, task)
}

// GetUserByID godoc
//	@Summary		Get task by ID
//	@Description	Retrieve a specific task by their unique identifier
//	@Tags			Tasks
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"User ID"
//	@Success		200	{object}	models.Task			"Successful response"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		404	{object}	response.Response	"User not found"
//	@Router			/tasks/{id} [get]
func (h *TaskHandlerImplementation) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, task)
}

// GetTasksByProject godoc
//	@Summary		Retrieve tasks by project ID
//	@Description	Retrieve paginated tasks associated with a specific project
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			project_id	path		int						true	"Project ID"
//	@Param			page		query		int						false	"Page number (default: 1)"
//	@Param			page_size	query		int						false	"Page size (default: 10)"
//	@Success		200			{object}	map[string]interface{}	"Paginated list of tasks"
//	@Failure		400			{object}	response.Response		"Invalid input"
//	@Failure		500			{object}	response.Response		"Server error"
//	@Router			/projects/{project_id}/tasks [get]
func (h *TaskHandlerImplementation) GetTasksByProject(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	projectID, _ := strconv.ParseUint(r.PathValue("project_id"), 10, 64)
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	tasks, total, err := h.service.GetTasksByProject(r.Context(), uint(projectID), page, pageSize)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"tasks": tasks,
		"total": total,
		"page":  page,
	})
}

// UpdateTask godoc
//	@Summary		Update a task
//	@Description	Update a task's details
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			task	body		models.Task			true	"Task details to be updated"
//	@Success		200		{object}	models.Task			"Task updated successfully"
//	@Failure		400		{object}	response.Response	"Invalid input"
//	@Failure		500		{object}	response.Response	"Server error"
//	@Router			/api/v1/tasks [put]
func (h *TaskHandlerImplementation) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.UpdateTask(r.Context(), &task); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, task)
}

// DeleteUser godoc
//	@Summary		Delete a task
//	@Description	Delete a task by its ID
//	@Tags			Tasks
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Task ID"
//	@Success		200	{object}	map[string]string	"Task deleted successfully"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		500	{object}	response.Response	"Server error"
//	@Router			/tasks/{id} [delete]
func (h *TaskHandlerImplementation) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.DeleteTask(r.Context(), uint(id)); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "task deleted successfully"})
}
