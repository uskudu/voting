package poll

import (
	"fmt"

	"gorm.io/gorm"
)

type RepositoryIface interface {
	CreatePoll(poll *Poll) error
	GetPolls() ([]Poll, error)
	GetPollByID(id string) (Poll, error)
	UpdatePoll(poll *Poll) error
	DeletePoll(id string) error
}

type Repository struct {
	DB *gorm.DB
}

func NewPollRepository(db *gorm.DB) RepositoryIface {
	return &Repository{DB: db}
}

func (r *Repository) CreatePoll(poll *Poll) error {
	return r.DB.Create(&poll).Error
}

func (r *Repository) GetPolls() ([]Poll, error) {
	var polls []Poll
	err := r.DB.Preload("Options").Find(&polls).Error
	return polls, err
}

func (r *Repository) GetPollByID(id string) (Poll, error) {
	var poll Poll
	err := r.DB.Preload("Options").First(&poll, "id = ?", id).Error
	return poll, err
}

func (r *Repository) UpdatePoll(poll *Poll) error {
	if err := r.DB.Where("poll_id = ?", poll.ID).Delete(&Option{}).Error; err != nil {
		return err
	}
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&poll).Error
}

func (r *Repository) DeletePoll(id string) error {
	result := r.DB.Delete(&Poll{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("poll not found")
	}
	return nil
}
