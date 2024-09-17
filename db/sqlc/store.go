package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries         // Embedding Queries to directly access its methods
	db       *sql.DB // The database connection pool
}

// NewStore creates a new Store object.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
// It accepts a context (ctx) and a function (fn) that takes *Queries and returns an error.
// If the function fn is executed successfully, the transaction is committed. If an error occurs, the transaction is rolled back.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Begin a new transaction using the provided context (ctx).
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		// Return the error if the transaction could not be started
		return err
	}

	// Create a new Queries object using the transaction (tx).
	// This ensures all queries in the fn function are executed within the transaction.
	q := New(tx)

	// Execute the passed function (fn), passing the transaction-bound Queries object (q).
	err = fn(q)
	if err != nil {
		// If an error occurs while executing the function, rollback the transaction.
		// If the rollback also fails, return both errors wrapped together.
		if rbErr := tx.Rollback(); rbErr != nil {
			// Combining the original error (err) with the rollback error (rbErr) for better error reporting.
			return fmt.Errorf("tx err: %v, rbErr: %v", err, rbErr)
		}
		// Return the original error if rollback was successful.
		return err
	}

	// If the function (fn) executes successfully, commit the transaction.
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `่json:"from_acc_id"`
	ToAccountId   int64 `่json:"to_acc_id"`
	Amount        int64 `่json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `่json:"transger"`
	FromAccount Account  `่json:"from_acc"`
	ToAccount   Account  `่json:"to_acc"`
	FromEntry   Entry    `่json:"from_entry"`
	ToEntry     Entry    `่json:"to_entry"`
}

var txKey = struct{}{}

// perform money transfer from 1 acc to another
// create transfer record, add account entries and update account balance in 1 transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccID: arg.FromAccountId,
			ToAccID:   arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccID:  arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccID:  arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		//update acc balance
		fmt.Println(txName, "update acc1 balance")

		if arg.FromAccountId < arg.ToAccountId { // handle case transfer each other
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId2,
		Amount: amount2,
	})
	return
}
