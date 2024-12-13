package services

import (
	"context"
	"example/project-management-system/internal/models"
	"example/project-management-system/internal/repositories"
	"fmt"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByID(ctx context.Context, id uint) (*models.Team, error)
	GetPaginatedTeams(ctx context.Context, page, pageSize int) ([]models.Team, int64, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uint) error
}

type TeamServiceImplementation struct {
	repo repositories.TeamRepository
}

func NewTeamService(repo repositories.TeamRepository) TeamService {
	return &TeamServiceImplementation{repo: repo}
}

func (s *TeamServiceImplementation) CreateTeam(ctx context.Context, team *models.Team) error {
	if team.Name == "" {
		return fmt.Errorf("team name is required")
	}
	return s.repo.CreateTeam(ctx, team)
}

func (s *TeamServiceImplementation) GetTeamByID(ctx context.Context, id uint) (*models.Team, error) {
	team, err := s.repo.GetTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("team not found")
	}
	return team, nil
}

func (s *TeamServiceImplementation) GetPaginatedTeams(ctx context.Context, page, pageSize int) ([]models.Team, int64, error) {
	return s.repo.GetAllTeams(ctx, page, pageSize)
}

func (s *TeamServiceImplementation) UpdateTeam(ctx context.Context, team *models.Team) error {
	if team.Name == "" {
		return fmt.Errorf("team name is required")
	}
	return s.repo.UpdateTeam(ctx, team)
}

func (s *TeamServiceImplementation) DeleteTeam(ctx context.Context, id uint) error {
	return s.repo.DeleteTeam(ctx, id)
}
