package migrations

import (
	"example/project-management-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func MigrateV1(tx *gorm.DB) error {
    err := tx.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
		&models.Team{},
	)
    if err != nil {
        return fmt.Errorf("v1 migration failed: %v", err)
    }

    return nil
}
