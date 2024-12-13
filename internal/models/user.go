package models

// internal/models/user.go

// User represents a user in the system
// @Description User model with basic information and relationships
type User struct {
	BaseModel
	// @Description Unique username for the user
	Username    string `json:"username" gorm:"unique;not null"`
	// @Description Unique email address of the user
	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"-" gorm:"not null"`
	// @Description User's first name
	FirstName   string `json:"first_name" `
	// @Description User's last name
	LastName    string `json:"last_name" `

	// Relationships are typically not serialized in Swagger docs
	ProjectIDs 	[]uint    `json:"project_ids" gorm:"-"`
	Projects    []Project `json:"projects" gorm:"many2many:user_projects;"`
	// Tasks       []Task    `gorm:"foreignKey:AssignedToID"`
	// Teams       []Team    `gorm:"many2many:user_teams;"`

	// @Description User's role in the system
	Role        string `json:"role" `
}
