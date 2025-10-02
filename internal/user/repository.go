package user

import (
	"fmt"

	"gorm.io/gorm"
)

type RepositoryIface interface {
	CreateUser(user *User) error
	GetUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	GetUserByUsername(username string) (User, error)
	UpdateUser(id, newUsername string) error
	DeleteUser(id string) error
}

type Repository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) RepositoryIface {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(user *User) error {
	return r.DB.Create(&user).Error
}

func (r *Repository) GetUsers() ([]User, error) {
	var users []User
	err := r.DB.Preload("Polls.Options").Find(&users).Error
	return users, err
}

func (r *Repository) GetUserByID(id string) (User, error) {
	var got User
	err := r.DB.First(&got, "id = ?", id).Error
	return got, err
}

func (r *Repository) GetUserByUsername(username string) (User, error) {
	var got User
	err := r.DB.First(&got, "username = ?", username).Error
	return got, err
}

func (r *Repository) UpdateUser(id, newUsername string) error {
	result := r.DB.Model(&User{}).
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
	result := r.DB.Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
