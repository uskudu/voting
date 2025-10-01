package user

import (
	"errors"

	"github.com/google/uuid"
)

type ServiceIface interface {
	CreateUser(username string) error
	GetUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	UpdateUser(id, newUsername string) error
	DeleteUser(id string) error
	Authenticate(id string) (*User, error)
}

type Service struct {
	repo RepositoryIface
}

func NewUserService(r RepositoryIface) ServiceIface {
	return &Service{repo: r}
}

func (s *Service) CreateUser(username string) error {
	newUser := User{
		ID:       uuid.NewString(),
		Username: username,
	}
	if err := s.repo.CreateUser(&newUser); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUsers() ([]User, error) {
	return s.repo.GetUsers()
}

func (s *Service) GetUserByID(id string) (User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) UpdateUser(id, newUsername string) error {
	return s.repo.UpdateUser(id, newUsername)
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *Service) Authenticate(username string) (*User, error) {
	got, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username")
	}
	return &got, nil
}
