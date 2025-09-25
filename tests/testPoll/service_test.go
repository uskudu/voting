package testPoll

import (
	"testing"
	"voting/internal/poll"

	"github.com/stretchr/testify/require"
)

func setupService(t *testing.T) poll.ServiceIface {
	repo := setupRepo(t)
	return poll.NewPollService(repo)
}

func TestCreatePollService(t *testing.T) {
	service := setupService(t)

	title := "test create poll"
	options := []poll.Option{
		poll.Option{
			Text: "a",
		},
		poll.Option{
			Text: "b",
		},
	}

	err := service.CreatePoll(title, options)
	require.NoError(t, err)

	//todo get created instance from db
	//rep := setupRepo(t)
	//var got poll.Poll
	//err = rep.(*poll.PollRepository).DB.Preload("Options").First(&got, "1").Error
	//require.NoError(t, err)
	//require.Equal(t, "test create poll", got.Title)
	//require.Len(t, got.Options, 2)
}

func TestGetPollsService(t *testing.T) {
	service := setupService(t)

	// todo create it without using service.CreatePoll
	title1 := "test get polls 1"
	options1 := []poll.Option{
		poll.Option{
			Text: "a",
		},
		poll.Option{
			Text: "b",
		},
	}
	title2 := "test get polls 2"
	options2 := []poll.Option{
		poll.Option{
			Text: "c",
		},
		poll.Option{
			Text: "d",
		},
	}
	err := service.CreatePoll(title1, options1)
	require.NoError(t, err)
	err = service.CreatePoll(title2, options2)
	require.NoError(t, err)

	polls, err := service.GetPolls()

	require.NoError(t, err)
	require.Len(t, polls, 2)
	require.Equal(t, "test get polls 1", polls[0].Title)
	require.Equal(t, "test get polls 2", polls[1].Title)
}

func TestGetPollByID(t *testing.T) {
	service := setupService(t)

	// todo create it without using service.CreatePoll
	title := "test get poll by id"
	options := []poll.Option{
		poll.Option{
			Text: "a",
		},
		poll.Option{
			Text: "b",
		},
	}
	err := service.CreatePoll(title, options)
	require.NoError(t, err)

	poll, err := service.GetPollByID("1")
	require.NoError(t, err)
	require.Equal(t, 1, poll.ID)
	require.Equal(t, "test get poll by id", poll.Title)
	require.Len(t, poll.Options, 2, "poll should have 2 options")

}
