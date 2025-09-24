package testPoll

import (
	"testing"
	"voting/internal/poll"
	"voting/tests/dbMock"

	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) poll.RepositoryIface {
	db := dbMock.SqliteMock()

	if err := db.AutoMigrate(&poll.Poll{}, &poll.Option{}); err != nil {
		panic(err)
	}

	return poll.NewPollRepository(db)
}

func TestCreatePoll(t *testing.T) {
	repo := setupRepo(t)

	p := &poll.Poll{
		Title: "test create poll",
		Options: []poll.Option{
			{Text: "a"},
			{Text: "b"},
		},
	}

	err := repo.CreatePoll(p)
	require.NoError(t, err)

	require.NotZero(t, p.ID, "poll should have ID after creation")

	var got poll.Poll
	err = repo.(*poll.PollRepository).DB.Preload("Options").First(&got, p.ID).Error

	require.NoError(t, err)
	require.Equal(t, "test create poll", got.Title)
	require.Len(t, got.Options, 2)
}

func TestGetPolls(t *testing.T) {
	repo := setupRepo(t)

	pollsToCreate := []poll.Poll{
		{
			Title: "test get polls 1",
			Options: []poll.Option{
				{Text: "a"},
				{Text: "b"},
			},
		},
		{
			Title: "test get polls 2",
			Options: []poll.Option{
				{Text: "c"},
				{Text: "d"},
			},
		},
	}

	for i := range pollsToCreate {
		err := repo.(*poll.PollRepository).DB.Create(&pollsToCreate[i]).Error
		require.NoError(t, err)
	}

	polls, err := repo.GetPolls()
	require.NoError(t, err)
	require.Len(t, polls, 2, "should return all polls")
	require.Equal(t, "test get polls 1", polls[0].Title)
	require.Equal(t, "test get polls 2", polls[1].Title)
}
