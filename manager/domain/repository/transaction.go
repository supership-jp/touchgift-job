//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package repository

import "context"

type TransactionHandler interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}
