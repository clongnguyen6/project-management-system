package handlers

import (
	"encoding/json"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/services"
	"example/project-management-system/internal/utils/response"
	"net/http"

	"strconv"
)

type TeamHandler interface {
	CreateTeam(w http.ResponseWriter, r *http.Request)
	GetTeamByID(w http.ResponseWriter, r *http.Request)
	GetPaginatedTeams(w http.ResponseWriter, r *http.Request)
	UpdateTeam(w http.ResponseWriter, r *http.Request)
	DeleteTeam(w http.ResponseWriter, r *http.Request)
}

type TeamHandlerImplementation struct {
	service services.TeamService
}

func NewTeamHandler(service services.TeamService) *TeamHandlerImplementation {
	return &TeamHandlerImplementation{service: service}
}

// CreateTeam godoc
//	@Summary		Create a new team
//	@Description	Create a new team with provided details
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			team	body		models.Team			true	"Team Creation Request"
//	@Success		201		{object}	map[string]int		"Team created successfully"
//	@Failure		400		{object}	response.Response	"Invalid input"
//	@Failure		500		{object}	response.Response	"Server error"
//	@Router			/teams [post]
func (h *TeamHandlerImplementation) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.CreateTeam(r.Context(), &team); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusCreated, team)
}

// func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
//     var input struct {
//         Name        string `json:"name"`
//         Description string `json:"description"`
//         ProjectID   uint   `json:"project_id"`
//         UserIDs     []uint `json:"user_ids"`
//     }

//     if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
//         response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
//         return
//     }

//     var users []models.User
//     if len(input.UserIDs) > 0 {
//         if err := db.Where("id IN ?", input.UserIDs).Find(&users).Error; err != nil {
//             response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
//             return
//         }
//     }

//     team := models.Team{
//         Name:        input.Name,
//         Description: input.Description,
//         ProjectID:   input.ProjectID,
//         Users:       users,
//     }

//     if err := h.service.CreateTeam(r.Context(), &team); err != nil {
//         response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
//         return
//     }

//     response.WriteJson(w, http.StatusCreated, team)
// }


// GetTeamByID godoc
//	@Summary		Get team by ID
//	@Description	Retrieve a specific team by their unique identifier
//	@Tags			Teams
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Team ID"
//	@Success		200	{object}	models.Team			"Successful response"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		404	{object}	response.Response	"User not found"
//	@Router			/teams/{id} [get]
func (h *TeamHandlerImplementation) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	team, err := h.service.GetTeamByID(r.Context(), uint(id))
	if err != nil {
		response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, team)
}


// GetPaginatedTeams godoc
//	@Summary		Get all teams
//	@Description	Retrieve paginated list of teams
//	@Tags			Teams
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int						false	"Page number"				default(1)
//	@Param			pageSize	query		int						false	"Number of teams per page"	default(10)
//	@Success		200			{object}	map[string]interface{}	"Successful response"
//	@Failure		400			{object}	response.Response		"Bad request"
//	@Router			/teams [get]
func (h *TeamHandlerImplementation) GetPaginatedTeams(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	teams, total, err := h.service.GetPaginatedTeams(r.Context(), page, pageSize)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]interface{}{
		"teams": teams,
		"total": total,
		"page":  page,
	})
}


// UpdateTeam godoc
//	@Summary		Update an existing team
//	@Description	Update a team's details by its ID
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Team ID"
//	@Param			team	body		models.Team			true	"Updated Team Details"
//	@Success		200		{object}	models.Team			"Successful response"
//	@Failure		400		{object}	response.Response	"Bad request"
//	@Failure		404		{object}	response.Response	"User not found"
//	@Router			/teams/{id} [put]
func (h *TeamHandlerImplementation) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.UpdateTeam(r.Context(), &team); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, team)
}


// DeleteTeam godoc
//	@Summary		Delete a team
//	@Description	Remove a team from the system by their ID
//	@Tags			Teams
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int					true	"Team ID to delete"
//	@Success		200	{object}	map[string]string	"Team deleted successfully"
//	@Failure		400	{object}	response.Response	"Bad request"
//	@Failure		500	{object}	response.Response	"Server error"
//	@Router			/teams/{id} [delete]
func (h *TeamHandlerImplementation) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := h.service.DeleteTeam(r.Context(), uint(id)); err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]string{"message": "team deleted successfully"})
}
