package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simple-bank/util"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))

	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmptyf(t, user, "User Cannot be empty")

	require.Equal(t, arg.Username, user.Username, "Username Not Equal")
	require.Equal(t, arg.HashedPassword, user.HashedPassword, "Hashed Password Not Equal")
	require.Equal(t, arg.FullName, user.FullName, "FullName Not Equal")
	require.Equal(t, arg.Email, user.Email, "Email Not Equal")

	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(),
		user1.Username)

	require.NoError(t, err)
	require.NotEmptyf(t, user2, "User 2 cannot be empty")

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.WithinDuration(t, user2.CreatedAt.Time, user1.CreatedAt.Time, time.Second, "Created at not equal")

}
