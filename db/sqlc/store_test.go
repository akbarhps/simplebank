package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		assert.Nil(t, err)

		result := <-results
		assert.NotEmpty(t, result)

		transfer := result.Transfer
		assert.Equal(t, acc1.ID, transfer.FromAccountID)
		assert.Equal(t, acc2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		assert.Nil(t, err)

		fromEntry := result.FromEntry
		assert.NotZero(t, fromEntry.ID)
		assert.Equal(t, acc1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		assert.Nil(t, err)

		toEntry := result.ToEntry
		assert.NotZero(t, toEntry.ID)
		assert.Equal(t, acc2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		assert.Nil(t, err)

		fromAccount := result.FromAccount
		assert.NotZero(t, fromAccount.ID)
		assert.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		assert.NotZero(t, toAccount.ID)
		assert.Equal(t, acc2.ID, toAccount.ID)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		// 5 times transaction should divisible by the amount
		// amount, 2 * amount, ....
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	assert.Nil(t, err)
	assert.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	assert.Nil(t, err)
	assert.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i%2 == 0 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.Nil(t, err)
	}

	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	assert.Nil(t, err)
	assert.Equal(t, acc1.Balance, updatedAcc1.Balance)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	assert.Nil(t, err)
	assert.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
