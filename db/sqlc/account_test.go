package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simple-bank/util"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmptyf(t, account, "Account Cannot be empty")

	require.Equal(t, arg.Owner, account.Owner, "Owner Not Equal")
	require.Equal(t, arg.Balance, account.Balance, "Balance Not Equal")
	require.Equal(t, arg.Currency, account.Currency, "Currency Not Equal")

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccountForOwner(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := GetAccountForOwnerParams{
		ID:    account1.ID,
		Owner: account1.Owner,
	}
	account2, err := testQueries.GetAccountForOwner(context.Background(),
		arg)

	require.NoError(t, err)
	require.NotEmptyf(t, account2, "Account 2 cannot be empty")

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second, "Created at not equal")
}
func TestQueries_GetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(),
		account1.ID)

	require.NoError(t, err)
	require.NotEmptyf(t, account2, "Account 2 cannot be empty")

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second, "Created at not equal")

}

func TestQueries_UpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.UpdateAccount(context.Background(),
		UpdateAccountParams{
			ID:      account1.ID,
			Balance: account1.Balance,
		})

	require.NoError(t, err)
	require.NotEmptyf(t, account2, "Account 2 cannot be empty")

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second, "Created at not equal")

}

func TestQueries_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)

	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, account2, "Account 2 expected to be empty")

}

func TestQueries_ListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestQueries_ListAccountsForOwner(t *testing.T) {
	var lastAcct Account
	for i := 0; i < 10; i++ {
		lastAcct = createRandomAccount(t)
	}

	arg := ListAccountsForOwnerParams{
		Owner:  lastAcct.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccountsForOwner(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastAcct.Owner)
	}
}
