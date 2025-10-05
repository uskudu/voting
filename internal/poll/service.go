package poll

import (
	"fmt"
)

type ServiceIface interface {
	CreatePoll(poll *Poll) error
	GetPolls() ([]Poll, error)
	GetPollByID(id string) (Poll, error)
	UpdatePoll(id string, poll Poll) error
	DeletePoll(id string) error
	AddVote(pollID, optionID string) error
}

type pollService struct {
	repo RepositoryIface
}

func NewPollService(r RepositoryIface) ServiceIface {
	return &pollService{repo: r}
}

func (s *pollService) CreatePoll(poll *Poll) error {
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

func (s *pollService) UpdatePoll(id string, poll Poll) error {
	pollFromDB, err := s.repo.GetPollByID(id)
	if err != nil {
		return err
	}
	for _, option := range pollFromDB.Options {
		if option.Votes != 0 {
			return fmt.Errorf("you cant edit poll containing one or more votes")
		}
	}

	pollFromDB.Title = poll.Title
	pollFromDB.Options = poll.Options

	if err = s.repo.UpdatePoll(&pollFromDB); err != nil {
		return err
	}
	return nil
}

func (s *pollService) DeletePoll(id string) error {
	return s.repo.DeletePoll(id)
}

func (s *pollService) AddVote(pollID, optionID string) error {
	return s.repo.AddVote(pollID, optionID)
}
