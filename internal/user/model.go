package user

import (
	"voting/internal/poll"
)

type User struct {
	ID       string      `gorm:"primaryKey" json:"id"`
	Username string      `gorm:"unique;not null" json:"username"`
	Polls    []poll.Poll `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"polls"`
}

type CreateOrLoginUserRequest struct {
	Username string `json:"username" binding:"required"`
}

type PatchUserRequest struct {
	Username string `json:"username"`
}
