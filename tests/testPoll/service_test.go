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

	got, err := service.GetPollByID("1")
	require.NoError(t, err)
	require.Equal(t, 1, got.ID)
	require.Equal(t, "test get poll by id", got.Title)
	require.Len(t, got.Options, 2, "poll should have 2 options")
}

func TestUpdatePoll(t *testing.T) {
	service := setupService(t)

	// todo create it without using service.CreatePoll
	title := "old title"
	options := []poll.Option{
		poll.Option{
			Text: "old option a",
		},
		poll.Option{
			Text: "old option b",
		},
	}
	err := service.CreatePoll(title, options)
	require.NoError(t, err)

	upd_poll := poll.Poll{
		Title: "updated title",
		Options: []poll.Option{
			{
				Text: "updated option a",
			},
			{
				Text: "updated option b",
			},
			{
				Text: "new option c",
			},
		},
	}

	err = service.UpdatePoll("1", upd_poll)
	require.NoError(t, err)

	got, err := service.GetPollByID("1")
	require.NoError(t, err)
	require.Equal(t, "updated title", got.Title)
	require.Len(t, got.Options, 3, "updated poll has 3 options")
	require.True(t, got.Options[0].Text == "updated option a")
	require.True(t, got.Options[1].Text == "updated option b")
	require.True(t, got.Options[2].Text == "new option c")
}
