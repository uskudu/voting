package poll

import (
	"fmt"
)

type ServiceIface interface {
	CreatePoll(title string, options []Option) error
	GetPolls() ([]Poll, error)
	GetPollByID(id string) (Poll, error)
	UpdatePoll(id string, poll Poll) error
	DeletePoll(id string) error
}

type pollService struct {
	repo RepositoryIface
}

func NewPollService(r RepositoryIface) ServiceIface {
	return &pollService{repo: r}
}

func (s *pollService) CreatePoll(title string, options []Option) error {
	poll := Poll{
		Title: title,
	}
	for _, opt := range options {
		option := Option{
			Text:   opt.Text,
			PollID: poll.ID,
		}
		poll.Options = append(poll.Options, option)
	}
	if err := s.repo.CreatePoll(&poll); err != nil {
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
			return fmt.Errorf("error while updating: you cant edit testPoll containing one or more votes")
		}
	}

	pollFromDB.Title = poll.Title
	pollFromDB.Options = poll.Options
	//for _, o := range testPoll.Options {
	//	pollFromDB.Options = append(pollFromDB.Options, Option{
	//		Text:   o.Text,
	//		PollID: pollFromDB.ID,
	//	})
	//}

	if err := s.repo.UpdatePoll(&pollFromDB); err != nil {
		return err
	}
	return nil
}

func (s *pollService) DeletePoll(id string) error {
	return s.repo.DeletePoll(id)
}
