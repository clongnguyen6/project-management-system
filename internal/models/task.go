package models

type Task struct {
	BaseModel
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ProjectID   uint    `json:"project_id"` // Many-to-One với Project
	Project     Project `gorm:"foreignKey:ProjectID" json:"project"`
	AssignedTo  uint    `json:"assigned_to"` // Many-to-One với User
	Assignee    User    `gorm:"foreignKey:AssignedTo" json:"assignee"`
}
