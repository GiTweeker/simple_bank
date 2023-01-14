package db

import (
	"context"
	"database/sql"
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

func TestUpdateUserOnlyFullName(t *testing.T) {
	olderUser := CreateRandomUser(t)
	newFullName := util.RandomOwner()

	arg := UpdateUserParams{
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Username: olderUser.Username,
	}
	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, olderUser.FullName, user.FullName)
	require.Equal(t, newFullName, user.FullName)
	require.Equal(t, olderUser.Email, user.Email)
	require.Equal(t, olderUser.HashedPassword, user.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	olderUser := CreateRandomUser(t)
	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: olderUser.Username,
	}
	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, olderUser.Email, user.Email)
	require.Equal(t, newEmail, user.Email)
	require.Equal(t, olderUser.FullName, user.FullName)
	require.Equal(t, olderUser.HashedPassword, user.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	olderUser := CreateRandomUser(t)
	newPassword := util.RandomString(6)
	newPasswordHash, err := util.HashPassword(newPassword)

	require.NoError(t, err)

	arg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newPasswordHash,
			Valid:  true,
		},
		Username: olderUser.Username,
	}
	user, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, olderUser.HashedPassword, user.HashedPassword)
	require.Equal(t, newPasswordHash, user.HashedPassword)
	require.Equal(t, olderUser.FullName, user.FullName)
	require.Equal(t, olderUser.Email, user.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	olderUser := CreateRandomUser(t)
	newPassword := util.RandomString(6)
	newPasswordHash, err := util.HashPassword(newPassword)
	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()

	require.NoError(t, err)

	arg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newPasswordHash,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Username: olderUser.Username,
	}
	user, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, olderUser.HashedPassword, user.HashedPassword)
	require.NotEqual(t, olderUser.FullName, user.FullName)
	require.NotEqual(t, olderUser.Email, user.Email)

	require.Equal(t, newPasswordHash, user.HashedPassword)
	require.Equal(t, newFullName, user.FullName)
	require.Equal(t, newEmail, user.Email)
}
