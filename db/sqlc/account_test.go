package db

import (
	"context"
	"database/sql"
	"github.com/akbarhps/simplebank/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) *Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotEmpty(t, acc)

	assert.Equal(t, acc.Owner, arg.Owner)
	assert.Equal(t, acc.Balance, arg.Balance)
	assert.Equal(t, acc.Currency, arg.Currency)
	assert.NotZero(t, acc.ID)
	assert.NotZero(t, acc.CreatedAt)

	return &acc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)

	getAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, getAcc)

	assert.Equal(t, acc.ID, getAcc.ID)
	assert.Equal(t, acc.Owner, getAcc.Owner)
	assert.Equal(t, acc.Balance, getAcc.Balance)
	assert.Equal(t, acc.Currency, getAcc.Currency)
	assert.WithinDuration(t, acc.CreatedAt, getAcc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      acc.ID,
		Balance: util.RandomMoney(),
	}

	updateAcc, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotEmpty(t, updateAcc)

	assert.Equal(t, acc.ID, updateAcc.ID)
	assert.Equal(t, acc.Owner, updateAcc.Owner)
	assert.Equal(t, arg.Balance, updateAcc.Balance)
	assert.Equal(t, acc.Currency, updateAcc.Currency)
	assert.WithinDuration(t, acc.CreatedAt, updateAcc.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc.ID)
	assert.Nil(t, err)

	getAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, getAcc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	listAcc, err := testQueries.ListAccounts(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotEmpty(t, listAcc)

	assert.Equal(t, int(arg.Limit), len(listAcc))

	for _, acc := range listAcc {
		assert.NotEmpty(t, acc)
	}
}
