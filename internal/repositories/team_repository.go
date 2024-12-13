package repositories

import (
	"context"
	"example/project-management-system/internal/models"

	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByID(ctx context.Context, id uint) (*models.Team, error)
	GetAllTeams(ctx context.Context, page, pageSize int) ([]models.Team, int64, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uint) error
}

type TeamRepositoryImplementation struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &TeamRepositoryImplementation{db: db}
}

func (r *TeamRepositoryImplementation) CreateTeam(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *TeamRepositoryImplementation) GetTeamByID(ctx context.Context, id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).
		Preload("Users").
		Preload("Project").
		First(&team, id).Error
	return &team, err
}

func (r *TeamRepositoryImplementation) GetAllTeams(ctx context.Context, page, pageSize int) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Team{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated records
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Preload("Users").
		Preload("Project").
		Find(&teams).Error

	return teams, total, err
}

func (r *TeamRepositoryImplementation) UpdateTeam(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *TeamRepositoryImplementation) DeleteTeam(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Team{}, id).Error
}
