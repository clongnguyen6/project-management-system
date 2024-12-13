package handlers

import (
	"encoding/json"
	"errors"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/services"
	"example/project-management-system/internal/utils/response"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
)

// UserHandler defines the interface for user-related HTTP handlers.
//	@title			UserHandler Interface
//	@description	Interface for handling user-related HTTP requests.
type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

// UserHandlerImplementation handles HTTP requests for CRUD operations against the user model.
type UserHandlerImplementation struct {
	userService services.UserService
}

// NewUserHandler creates a new instance of the UserHandlerImplementation.
func NewUserHandler(userService services.UserService) *UserHandlerImplementation {
	return &UserHandlerImplementation{
		userService: userService,
	}
}

/* TODO:
//	@Param			email		formData	string				true	"Email"
//	@Param			password	formData	string				true	"Password"
//	@Param			firstname	formData	string				true	"Firstname"
//	@Param			lastname	formData	string				true	"Lastname"
//	@Param			project_ids	formData	[]int				true	"Project IDs"
//	@Param			role		formData	string				true	"Role"
*/
// CreateUser godoc
//	@Summary		Create a new user
//	@Description	Create a new user with provided details
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user	body		models.User			true	"Username"
//	@Success		201		{object}	map[string]int		"User created successfully"
//	@Failure		400		{object}	response.Response	"Invalid input"
//	@Failure		500		{object}	response.Response	"Server error"
//	@Router			/users [post]
func (s *UserHandlerImplementation) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if errors.Is(err, io.EOF) {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
		return
	}

	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := validator.New().Struct(user); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
		return
	}

	if err := s.userService.CreateUser(r.Context(), &user); err != nil {
		response.WriteJson(
			w, http.StatusInternalServerError, response.GeneralError(
				fmt.Errorf("%s", err.Error()),
			),
		)
		return
	}

	slog.Info("user created successfully", slog.String("userId", fmt.Sprint(user.ID)))

	response.WriteJson(w, http.StatusCreated, user)
}

// GetUsers godoc
//	@Summary		Get all users
//	@Description	Retrieve paginated list of users
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int						false	"Page number"				default(1)
//	@Param			pageSize	query		int						false	"Number of users per page"	default(10)
//	@Success		200			{object}	map[string]interface{}	"Successful response"
//	@Failure		400			{object}	response.Response		"Bad request"
//	@Router			/users [get]
func (s *UserHandlerImplementation) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteJson(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("%s", "Method not allowed")))
		return
	}

	// Parse query parameters for pagination
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	users, total, err := s.userService.GetAllUsers(r.Context(), page, pageSize)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("%s", "Failed to list users")))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
	})
}

// GetUserByID godoc
//	@Summary		Get user by ID
//	@Description	Retrieve a specific user by their unique identifier
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"User ID"
//	@Success		200	{object}	models.User			"Successful response"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		404	{object}	response.Response	"User not found"
//	@Router			/users/{id} [get]
func (s *UserHandlerImplementation) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("%s", "Invalid user ID")))
		return
	}

	slog.Info("getting a student", slog.String("id", strconv.FormatUint(id, 10)))

	if strconv.FormatUint(id, 10) == "" {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("id cannot empty")))
		return
	}

	user, err := s.userService.GetUserByID(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("could not get the user")))
		return
	}

	response.WriteJson(w, http.StatusOK, user)
}

// DeleteUser godoc
//	@Summary		Delete an user
//	@Description	Remove an user from the system by their ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"User ID to delete"
//	@Success		200	{object}	map[string]string	"User deleted successfully"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		500	{object}	response.Response	"Server error"
//	@Router			/users/{id} [delete]
func (s *UserHandlerImplementation) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("%s", "Invalid user ID")))
		return
	}

	slog.Info("deleting a user", slog.String("id", strconv.FormatUint(id, 10)))

	if strconv.FormatUint(id, 10) == "" {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("id cannot be empty")))
		return
	}

	err = s.userService.DeleteUser(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("could not delete user")))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "user delete successfully"})
}
