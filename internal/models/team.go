package models

import "time"

// Team Model (Many-to-Many with User, Many-to-One with Project)
type Team struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	ProjectID   uint	`json:"project_id"`
	Project     Project `json:"project" gorm:"foreignKey:ProjectID;constraint:onUpdate:CASCADE,onDelete:SET NULL;"` // onUpdate:CASCADE: Khi ProjectID trong bảng Project thay đổi, nó sẽ cập nhật tự động trong bảng Team. onDelete:SET NULL: Nếu một project bị xóa, ProjectID trong bảng Team sẽ được đặt thành NULL thay vì xóa toàn bộ team.
	Users       []User `json:"users" gorm:"many2many:user_teams;constraint:onUpdate:CASCADE,onDelete:CASCADE;"` // onUpdate:CASCADE: Khi UserID trong bảng User thay đổi, liên kết trong bảng trung gian (user_teams) sẽ được cập nhật. onDelete:CASCADE: Khi một user bị xóa, liên kết trong bảng trung gian (user_teams) cũng bị xóa.
}