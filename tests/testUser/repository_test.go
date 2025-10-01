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
	ID:       1,
	Username: "test user",
}

func TestCreateUserRepo(t *testing.T) {
	repo := setupRepo(t)

	err := repo.CreateUser(u)
	require.NoError(t, err)

	require.NotZero(t, u.ID, "user should have ID after creation")

	var got user.User
	err = repo.(*user.UserRepository).DB.First(&got, u.ID).Error

	require.NoError(t, err)
	require.Equal(t, "test user", got.Username)
}
func TestGetUsersRepo(t *testing.T) {
	repo := setupRepo(t)

	usersToCreate := []user.User{
		{Username: "alice"},
		{Username: "bob"},
	}

	for i := range usersToCreate {
		err := repo.CreateUser(&usersToCreate[i])
		require.NoError(t, err)
	}

	users, err := repo.GetUsers()
	require.NoError(t, err)
	require.Len(t, users, 2, "should return all users")
	require.Equal(t, "alice", users[0].Username)
	require.Equal(t, "bob", users[1].Username)
}

func TestGetUserByIDRepo(t *testing.T) {
	repo := setupRepo(t)

	err := repo.CreateUser(u)
	require.NoError(t, err)

	got, err := repo.GetUserByID("1")
	require.NoError(t, err)
	require.Equal(t, "test user", got.Username)
}

func TestUpdateUserRepo(t *testing.T) {
	repo := setupRepo(t)

	err := repo.CreateUser(u)
	require.NoError(t, err)

	old, err := repo.GetUserByID("1")
	require.NoError(t, err)
	require.Equal(t, "test user", old.Username)

	err = repo.UpdateUser("1", "updated user")
	require.NoError(t, err)

	updated, err := repo.GetUserByID("1")
	require.NoError(t, err)
	require.Equal(t, "updated user", updated.Username)
}

func TestDeleteUserRepo(t *testing.T) {
	repo := setupRepo(t)

	err := repo.CreateUser(u)
	require.NoError(t, err)

	err = repo.DeleteUser("1")
	require.NoError(t, err)

	got, err := repo.GetUserByID("1")
	require.Error(t, err)
	require.Empty(t, got)
}
