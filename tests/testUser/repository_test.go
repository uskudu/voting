package testUser

import (
	"testing"
	"voting/internal/user"
	"voting/tests/dbMock"

	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) user.RepositoryIface {
	db := dbMock.SqliteMock()
	if err := db.AutoMigrate(&user.User{}); err != nil {
		panic(err)
	}
	return user.NewUserRepository(db)
}

var u = &user.User{
	Username: "test user",
}

func TestCreatePollRepo(t *testing.T) {
	repo := setupRepo(t)

	err := repo.CreateUser(u)
	require.NoError(t, err)

	require.NotZero(t, u.ID, "poll should have ID after creation")

	var got user.User
	err = repo.(*user.UserRepository).DB.First(&got, u.ID).Error

	require.NoError(t, err)
	require.Equal(t, "test user", got.Username)
}

//
//func TestGetPollsRepo(t *testing.T) {
//	repo := setupRepo(t)
//
//	pollsToCreate := []poll.Poll{
//		{
//			Title: "test get polls 1",
//			Options: []poll.Option{
//				{Text: "a"},
//				{Text: "b"},
//			},
//		},
//		{
//			Title: "test get polls 2",
//			Options: []poll.Option{
//				{Text: "c"},
//				{Text: "d"},
//			},
//		},
//	}
//
//	for i := range pollsToCreate {
//		err := repo.(*poll.PollRepository).DB.Create(&pollsToCreate[i]).Error
//		require.NoError(t, err)
//	}
//
//	polls, err := repo.GetPolls()
//	require.NoError(t, err)
//	require.Len(t, polls, 2, "should return all polls")
//	require.Equal(t, "test get polls 1", polls[0].Title)
//	require.Equal(t, "test get polls 2", polls[1].Title)
//}
//
//func TestGetPollByIDRepo(t *testing.T) {
//	repo := setupRepo(t)
//
//	err := repo.CreatePoll(p)
//	require.NoError(t, err)
//
//	got, err := repo.GetPollByID("1")
//	require.NoError(t, err)
//	require.Equal(t, "test poll", got.Title)
//	require.Len(t, got.Options, 2, "should get all options")
//}
//
//func TestUpdatePollRepo(t *testing.T) {
//	repo := setupRepo(t)
//
//	err := repo.CreatePoll(p)
//	require.NoError(t, err)
//	// check if old is ok
//	old, err := repo.GetPollByID("1")
//	require.NoError(t, err)
//	require.Equal(t, "test poll", old.Title)
//	require.Len(t, old.Options, 2, "old has two options")
//
//	upd := &poll.Poll{
//		ID:    1,
//		Title: "updated title",
//		Options: []poll.Option{
//			{Text: "the only option"},
//		},
//	}
//	err = repo.UpdatePoll(upd)
//	require.NoError(t, err)
//	// check if updated is ok
//	updated, err := repo.GetPollByID("1")
//	require.NoError(t, err)
//	require.Equal(t, "updated title", updated.Title)
//	require.Len(t, updated.Options, 1, "updated has only one option")
//	require.True(t, updated.Options[0].Text == "the only option")
//}
//
//func TestDeletePollRepo(t *testing.T) {
//	repo := setupRepo(t)
//
//	err := repo.CreatePoll(p)
//	require.NoError(t, err)
//
//	err = repo.DeletePoll("1")
//	require.NoError(t, err)
//
//	got, err := repo.GetPollByID("1")
//	require.Error(t, err)
//	require.Empty(t, got)
//}

// todo somewhere id is int somewhere is string
