package testPoll

import (
	"testing"
	"voting/internal/poll"
	"voting/tests/dbMock"

	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) poll.RepositoryIface {
	db := dbMock.SqliteMock()
	return poll.NewPollRepository(db)
}

func TestCreatePoll(t *testing.T) {
	repo := setupRepo(t)

	p := poll.Poll{
		Title: "test poll",
		Options: []poll.Option{
			{Text: "a"},
			{Text: "b"},
		},
	}

	err := repo.CreatePoll(p)
	require.NoError(t, err)

	var got poll.Poll
	err = repo.(*poll.PollRepository).DB.First(&got, "id = ?", 1).Error
	require.NoError(t, err)

	require.Equal(t, "test poll", got.Title)
}
