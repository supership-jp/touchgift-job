package repository

import "context"

type TransactionHandler interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}