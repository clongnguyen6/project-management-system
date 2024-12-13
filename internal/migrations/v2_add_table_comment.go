package migrations

import (
	"example/project-management-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func MigrateV2(tx *gorm.DB) error {
    if !tx.Migrator().HasTable(&models.Comment{}) {
        err := tx.Migrator().CreateTable(&models.Comment{})
        if err != nil {
            return fmt.Errorf("v2 migration failed to create comments table")
        }
    }

    return nil
}
