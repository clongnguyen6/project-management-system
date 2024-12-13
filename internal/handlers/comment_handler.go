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

type CommentHandler interface {
	CreateComment(w http.ResponseWriter, r *http.Request)
	GetCommentByID(w http.ResponseWriter, r *http.Request)
	GetCommentsByTask(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
}

type CommentHandlerImplementation struct {
	service services.CommentService
}

func NewCommentHandler(service services.CommentService) *CommentHandlerImplementation {
	return &CommentHandlerImplementation{service: service}
}

// CreateComment godoc
//	@Summary		Create a new comment
//	@Description	Create a new comment associated with a specific task
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			comment	body		models.Comment		true	"Comment Creation Request"
//	@Success		201		{object}	models.Comment		"Comment created successfully"
//	@Failure		400		{object}	response.Response	"Invalid input"
//	@Failure		500		{object}	response.Response	"Server error"
//	@Router			/comments [post]
func (h *CommentHandlerImplementation) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.CreateComment(r.Context(), &comment); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusCreated, comment)
}

// GetCommentByID godoc
//	@Summary		Get a comment by ID
//	@Description	Retrieve a comment using its unique ID
//	@Tags			Comments
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Comment ID"
//	@Success		200	{object}	models.Comment		"Successful response"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		404	{object}	response.Response	"Comment not found"
//	@Router			/comments/{id} [get]
func (h *CommentHandlerImplementation) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("%s", "invalid ID")))
		return
	}

	comment, err := h.service.GetCommentByID(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, comment)
}

// GetCommentsByTask godoc
//	@Summary		Get comments by task
//	@Description	Retrieve paginated comments associated with a specific task
//	@Tags			Comments
//	@Produce		json
//	@Security		BearerAuth
//	@Param			task_id		path		int						true	"Task ID"
//	@Param			page		query		int						false	"Page number"					default(1)
//	@Param			page_size	query		int						false	"Number of comments per page"	default(10)
//	@Success		200			{object}	map[string]interface{}	"Successful response"
//	@Failure		400			{object}	response.Response		"Bad request"
//	@Router			/tasks/{task_id}/comments [get]
func (h *CommentHandlerImplementation) GetCommentsByTask(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	taskID, _ := strconv.ParseUint(r.PathValue("task_id"), 10, 64)
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	comments, total, err := h.service.GetCommentsByTask(r.Context(), uint(taskID), page, pageSize)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"comments": comments,
		"total":    total,
		"page":     page,
	})
}

// DeleteComment godoc
//	@Summary		Delete a comment
//	@Description	Remove a comment from the system by its ID
//	@Tags			Comments
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Comment ID to delete"
//	@Success		200	{object}	map[string]string	"Comment deleted successfully"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		500	{object}	response.Response	"Server error"
//	@Router			/comments/{id} [delete]
func (h *CommentHandlerImplementation) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.DeleteComment(r.Context(), uint(id)); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "comment deleted successfully"})
}
