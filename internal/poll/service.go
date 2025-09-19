package poll

import (
	"fmt"

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

func (s *pollService) GetPolls() ([]Poll, error) {
	return s.repo.GetPolls()
}

func (s *pollService) GetPollByID(id string) (Poll, error) {
	return s.repo.GetPollByID(id)
}

func (s *pollService) UpdatePoll(id, title string, options []string) error {
	pollFromDB, err := s.repo.GetPollByID(id)
	if err != nil {
		return err
	}
	for _, option := range pollFromDB.Options {
		if option.Votes != 0 {
			return fmt.Errorf("error while updating: you cant edit poll containing one or more votes")
		}
	}

	pollFromDB.Title = title
	pollFromDB.Options = []Option{}
	for _, o := range options {
		pollFromDB.Options = append(pollFromDB.Options, Option{
			ID:     uuid.NewString(),
			Text:   o,
			PollID: pollFromDB.ID,
		})
	}

	if err := s.repo.UpdatePoll(pollFromDB); err != nil {
		return err
	}
	return nil
}

func (s *pollService) DeletePoll(id string) error {
	return s.repo.DeletePoll(id)
}
