package models

import (
	"time"
)

// Comment Model (Many-to-One with Task, User)
type Comment struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Content   string     `json:"content" gorm:"not null"`
	TaskID    uint       `json:"task_id"`
	Task      Task       `json:"task" gorm:"foreignKey:TaskID;constraint:onUpdate:CASCADE,onDelete:CASCADE;"`
	UserID    uint       `json:"user_id"`
	User      User       `json:"user" gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:SET NULL;"`

}