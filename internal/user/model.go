package user

import "voting/internal/poll"

type User struct {
	ID       int         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string      `json:"username"`
	Polls    []poll.Poll `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"polls"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
}
