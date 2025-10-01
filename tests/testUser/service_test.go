package testUser

import (
	"testing"
	"voting/internal/user"

	"github.com/stretchr/testify/require"
)

func setupService(t *testing.T) user.ServiceIface {
	repo := setupRepo(t)
	return user.NewUserService(repo)
}

var username = "test user"

func TestCreateUserService(t *testing.T) {
	service := setupService(t)

	err := service.CreateUser("test user")
	require.NoError(t, err)

	got, err := service.GetUserByID("1")
	require.NoError(t, err)
	require.Equal(t, "test user", got.Username)
}

func TestGetUsersService(t *testing.T) {
	service := setupService(t)

	err := service.CreateUser("alice")
	require.NoError(t, err)
	err = service.CreateUser("bob")
	require.NoError(t, err)

	users, err := service.GetUsers()
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "alice", users[0].Username)
	require.Equal(t, "bob", users[1].Username)
}

func TestGetUserByIDService(t *testing.T) {
	service := setupService(t)

	err := service.CreateUser(username)
	require.NoError(t, err)

	got, err := service.GetUserByID("1")
	require.NoError(t, err)
	require.Equal(t, 1, got.ID)
	require.Equal(t, "test user", got.Username)
}

func TestUpdateUserService(t *testing.T) {
	service := setupService(t)

	err := service.CreateUser(username)
	require.NoError(t, err)

	err = service.UpdateUser("1", "updated user")
	require.NoError(t, err)

	got, err := service.GetUserByID("1")
	require.NoError(t, err)
	require.Equal(t, "updated user", got.Username)
}

func TestDeleteUserService(t *testing.T) {
	service := setupService(t)

	err := service.CreateUser(username)
	require.NoError(t, err)

	err = service.DeleteUser("1")
	require.NoError(t, err)

	got, err := service.GetUserByID("1")
	require.Error(t, err)
	require.Empty(t, got)
}
