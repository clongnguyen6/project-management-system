package models

import (
	"time"
)

// Project Model (Many-to-Many with User, One-to-Many with Task, Team)
type Project struct {
	BaseModel
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	UserIDs		[]uint	  `json:"user_ids" gorm:"-"`
	Users       []User    `json:"users" gorm:"many2many:user_projects;"`
	// One-to-Many with Tasks
	Tasks 		[]Task 	   `json:"tasks" gorm:"foreignKey:ProjectID"`
	// One-to-Many with Teams
	Teams       []Team    `json:"teams" gorm:"foreignKey:ProjectID"`
}
