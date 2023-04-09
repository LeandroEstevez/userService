package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
)

// Store gives all functions to execute db queries and transactions
type Store interface {
	Querier
	DeleteUserTx(ctx context.Context, username string) error
	UpdateTotalTx(ctx context.Context, entry Entry) error
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

type UpdatedEntryMessage struct {
	OrignalEntry Entry
	UpdatedEntry Entry
}

// Creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

func (store *SQLStore) DeleteUserTx(ctx context.Context, username string) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := q.DeleteUser(ctx, username)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Updates the amount of an user and updates the total expense in the user
func (store *SQLStore) UpdateTotalTx(ctx context.Context, entry Entry) error {
	err := store.execTx(ctx, func(q *Queries) error {
		user, err := store.GetUserForUpdate(ctx, entry.Owner)
		if err != nil {
			fmt.Println("There was an error", err)
			return err
			// TODO: send back an event to entriesMicroService with the error
		}

		totalExpenseUpdate := user.TotalExpenses + entry.Amount

		params := UpdateUserParams{
			Username:      entry.Owner,
			TotalExpenses: totalExpenseUpdate,
		}
		_, err = store.UpdateUser(ctx, params)
		if err != nil {
			fmt.Println("There was an error", err)
			return err
			// TODO: send back an event to entriesMicroService with the error
		}

		return nil
	})

	return err
}
