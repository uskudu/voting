package poll

import (
	"github.com/google/uuid"
)

type PollIface interface {
	CreatePoll(title string, options []string) error
	GetPolls() ([]Poll, error)
	GetPollByID(id string) (Poll, error)
	UpdatePoll(id, title string, options []string) error
	DeletePoll(id string) error
}

type pollService struct {
	repo RepositoryIface
}

func NewPollService(r RepositoryIface) PollIface {
	return &pollService{repo: r}
}

func (s *pollService) CreatePoll(title string, options []string) error {
	poll := Poll{
		ID:    uuid.NewString(),
		Title: title,
	}
	for _, text := range options {
		option := Option{
			ID:     uuid.NewString(),
			Text:   text,
			PollID: poll.ID,
		}
		poll.Options = append(poll.Options, option)
	}
	if err := s.repo.CreatePoll(poll); err != nil {
		return err
	}
	return nil
}

func GetPolls() ([]Poll, error) {

}
func GetPollByID(id string) (Poll, error) {

}
func UpdatePoll(id, title string, options []string) error {

}
func DeletePoll(id string) error {

}
