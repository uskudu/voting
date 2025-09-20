package poll

import (
	"fmt"

	"gorm.io/gorm"
)

type RepositoryIface interface {
	CreatePoll(poll Poll) error
	GetPolls() ([]Poll, error)
	GetPollByID(id string) (Poll, error)
	UpdatePoll(poll Poll) error
	DeletePoll(id string) error
}

type pollRepository struct {
	db *gorm.DB
}

func NewPollRepository(db *gorm.DB) RepositoryIface {
	return &pollRepository{db: db}
}

func (r *pollRepository) CreatePoll(poll Poll) error {
	return r.db.Create(&poll).Error
}

func (r *pollRepository) GetPolls() ([]Poll, error) {
	var polls []Poll
	err := r.db.Preload("Options").Find(&polls).Error
	return polls, err
}

func (r *pollRepository) GetPollByID(id string) (Poll, error) {
	var poll Poll
	err := r.db.Preload("Options").First(&poll, "id = ?", id).Error
	return poll, err
}

func (r *pollRepository) UpdatePoll(poll Poll) error {
	if err := r.db.Where("poll_id = ?", poll.ID).Delete(&Option{}).Error; err != nil {
		return err
	}
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&poll).Error
}

func (r *pollRepository) DeletePoll(id string) error {
	result := r.db.Delete(&Poll{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("poll not found")
	}
	return nil
}
