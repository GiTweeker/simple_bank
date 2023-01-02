package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err : %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct {
}{}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntryForAccount(ctx, CreateEntryForAccountParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntryForAccount(ctx, CreateEntryForAccountParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount,
				err = addMoneyToAccount(ctx, q, arg.FromAccountId, arg.ToAccountId, -arg.Amount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount,
				err = addMoneyToAccount(ctx, q, arg.ToAccountId, arg.FromAccountId, arg.Amount, -arg.Amount)
		}

		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

func addMoneyToAccount(ctx context.Context, q *Queries,
	accountId1, accountId2 int64, amount1, amount2 int64) (account1, account2 Account, err error) {
	if amount2+amount1 != 0 {
		return account1, account2, errors.New("amount 1 and Amount 2 must sum to 0")
	}
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountId2,
		Amount: +amount2,
	})
	if err != nil {
		return
	}

	return
}
