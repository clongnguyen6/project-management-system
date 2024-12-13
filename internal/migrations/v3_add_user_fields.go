package migrations

import (
	"example/project-management-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func MigrateV3(tx *gorm.DB) error {
    if !tx.Migrator().HasTable(&models.Task{}) {
        err := tx.Migrator().CreateTable(&models.Task{})
        if err != nil {
            return fmt.Errorf("v2 migration failed to create tasks table")
        }
    }

    if !tx.Migrator().HasColumn(&models.User{}, "Projects") {
        err := tx.Migrator().AddColumn(&models.User{}, "Projects")
        if err != nil {
            return fmt.Errorf("v2 migration failed to add projects column for users: %v", err)
        }
    }

    if !tx.Migrator().HasColumn(&models.Project{}, "Users") {
        err := tx.Migrator().AddColumn(&models.Project{}, "Users")
        if err != nil {
            return fmt.Errorf("v2 migration failed to add users column for projects: %v", err)
        }
    }

    if !tx.Migrator().HasColumn(&models.Project{}, "Tasks") {
        err := tx.Migrator().AddColumn(&models.Project{}, "Tasks")
        if err != nil {
            return fmt.Errorf("v2 migration failed to add tasks column for projects: %v", err)
        }
    }

    if !tx.Migrator().HasColumn(&models.Project{}, "Teams") {
        err := tx.Migrator().AddColumn(&models.Project{}, "Teams")
        if err != nil {
            return fmt.Errorf("v2 migration failed to add teams column for projects: %v", err)
        }
    }

    return nil
}
