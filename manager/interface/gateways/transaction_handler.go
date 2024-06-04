//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/$GOPACKAGE/$GOFILE
package gateways

import (
	"context"
	"touchgift-job-manager/domain/models/repository"
)

type TransactionHandler interface {
	Begin(ctx context.Context) (repository.Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}
