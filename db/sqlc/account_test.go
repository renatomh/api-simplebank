package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Function to be used during the tests
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	// Setting up data to be used during the account's creation
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// Creating the account item
	account, err := testQueries.CreateAccount(context.Background(), arg)
	// In order for the tests to pass, there should be no errors, it shouldn't be empty
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// And it should meet the required data
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// Cheking if ID and "created at" fields were filled
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	// Returning the created account
	return account
}

func TestCreateAccount(t *testing.T) {
	// Calling the function to create a random account and check if everything is ok
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	// First, we should create an account
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// In order for the tests to pass, there should be no errors, it shouldn't be empty
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// Checking if fields match
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	// First, we should create an account
	account1 := createRandomAccount(t)

	// Defining data to be updated
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	// In order for the tests to pass, there should be no errors, it shouldn't be empty
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// Checking if fields match
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	// First, we should create an account
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	// In order for the tests to pass, there should be no errors
	require.NoError(t, err)

	// Now, when querying the data, it should be empty
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	// First, we create random accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	// Defining pagination settings
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	// In order for the tests to pass, there should be no errors
	require.NoError(t, err)
	// And 5 items must be returned
	require.Len(t, accounts, 5)

	// Also, no account should be empty
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
