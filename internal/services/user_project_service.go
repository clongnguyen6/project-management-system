package services

import "example/project-management-system/internal/repositories"

type UserProjectService interface {
	AddUserToProject(userID uint, projectID uint) error
}

type UserProjectServiceImplementation struct {
	userProjectRepository repositories.UserProjectRepository
}

func NewUserProjectService(userProjectRepository repositories.UserProjectRepository) UserProjectService {
	return &UserProjectServiceImplementation{userProjectRepository: userProjectRepository}
}

// AddUserToProject adds a user to a project.
func (service *UserProjectServiceImplementation) AddUserToProject(userID uint, projectID uint) error {
	return service.userProjectRepository.AddUserToProject(userID, projectID)
}
