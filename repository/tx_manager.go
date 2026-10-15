package repository

import "context"

// TransactionManager executes a piece of work inside a database transaction.
type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, work func(ctx context.Context) error) error
}
