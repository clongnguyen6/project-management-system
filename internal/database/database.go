package database

import (
	"example/project-management-system/internal/config"
	"example/project-management-system/internal/migrations"
	"example/project-management-system/pkg/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
	Log logger.Logger
}

type MigrationVersion struct {
    Version int `gorm:"primaryKey"`
}

func NewPostgresConnection(cfg *config.Config, logger logger.Logger) *Database {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password, cfg.Db.DBName, cfg.Db.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect database", "error", err)
	}

	// Auto migrate models
	err = Migrate(db)

	if err != nil {
		logger.Fatal("Failed to migrate database", "error", err)
	}

	return &Database{
		DB:  db,
		Log: logger,
	}
}

func Migrate(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
        return fmt.Errorf("failed to begin transaction: %v", tx.Error)
    }

    // Rollback if err
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()


    err := tx.AutoMigrate(&MigrationVersion{})
    if err != nil {
		tx.Rollback()
        return err
    }

    migrationFuncs := []func(*gorm.DB) error{
        migrations.MigrateV1,
        migrations.MigrateV2,
    }

    for i, migrate := range migrationFuncs {
        version := i + 1
        
        var existingVersion MigrationVersion
        err := tx.Where("version = ?", version).First(&existingVersion).Error
        if err == gorm.ErrRecordNotFound {
            err := migrate(tx)
            if err != nil {
				tx.Rollback()
                return fmt.Errorf("migration v%d failed: %v", version, err)
            }

            tx.Create(&MigrationVersion{Version: version})
        }
    }

    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("databae failed to commit transaction: %v", err)
    }

    return nil
}
