package crud

import (
	"voting/internal/user"

	"github.com/google/uuid"
)

type ServiceIface interface {
	CreateUser(username string) error
	GetUsers() ([]user.User, error)
	GetUserByID(id string) (user.User, error)
	UpdateUser(id, newUsername string) error
	DeleteUser(id string) error
}

type userService struct {
	repo RepositoryIface
}

func NewUserService(r RepositoryIface) ServiceIface {
	return &userService{repo: r}
}

func (s *userService) CreateUser(username string) error {
	newUser := user.User{
		ID:       uuid.NewString(),
		Username: username,
	}
	if err := s.repo.CreateUser(&newUser); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUsers() ([]user.User, error) {
	return s.repo.GetUsers()
}

func (s *userService) GetUserByID(id string) (user.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) UpdateUser(id, newUsername string) error {
	return s.repo.UpdateUser(id, newUsername)
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
