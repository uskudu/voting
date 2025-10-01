package crud

import (
	"fmt"
	"voting/internal/user"

	"gorm.io/gorm"
)

type RepositoryIface interface {
	CreateUser(user *user.User) error
	GetUsers() ([]user.User, error)
	GetUserByID(id string) (user.User, error)
	UpdateUser(id, newUsername string) error
	DeleteUser(id string) error
}

type Repository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) RepositoryIface {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(user *user.User) error {
	return r.DB.Create(&user).Error
}

func (r *Repository) GetUsers() ([]user.User, error) {
	var users []user.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *Repository) GetUserByID(id string) (user.User, error) {
	var got user.User
	err := r.DB.First(&got, "id = ?", id).Error
	return got, err
}

func (r *Repository) UpdateUser(id, newUsername string) error {
	result := r.DB.Model(&user.User{}).
		Where("id = ?", id).
		Update("username", newUsername)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *Repository) DeleteUser(id string) error {
	result := r.DB.Delete(&user.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
