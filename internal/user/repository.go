package user

import (
	"fmt"

	"gorm.io/gorm"
)

type RepositoryIface interface {
	CreateUser(user *User) error
	GetUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	UpdateUser(id, newUsername string) error
	DeleteUser(id string) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) RepositoryIface {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *User) error {
	return r.DB.Create(&user).Error
}

func (r *UserRepository) GetUsers() ([]User, error) {
	var users []User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUserByID(id string) (User, error) {
	var user User
	err := r.DB.First(&user, "id = ?", id).Error
	return user, err
}

func (r *UserRepository) UpdateUser(id, newUsername string) error {
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

func (r *UserRepository) DeleteUser(id string) error {
	result := r.DB.Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
