// Package gormrepo provides a GORM-based repository toolkit: DB management,
// transactions, generic CRUD helpers, query scopes and option conversion.
package gormrepo

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/banbridge/common/repository"
)

// DBManager wraps a *gorm.DB and provides transaction and master/slave
// selection through context.
type DBManager struct {
	db *gorm.DB
}

// NewManager creates a DBManager from an existing *gorm.DB. The caller keeps
// ownership of the underlying database.
func NewManager(db *gorm.DB) (*DBManager, error) {
	if db == nil {
		return nil, fmt.Errorf("gormrepo: db is nil")
	}
	return &DBManager{db: db}, nil
}

type txKey struct{}

func (m *DBManager) inTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// UseTx attaches a transaction to the context so GetDB returns the transaction
// instead of the shared database.
func UseTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// ExecuteInTransaction executes work inside a transaction. The context passed
// to work is bound to the transaction.
func (m *DBManager) ExecuteInTransaction(ctx context.Context, work func(ctx context.Context) error) error {
	return m.GetDB(ctx).Transaction(func(itx *gorm.DB) error {
		return work(m.inTx(ctx, itx))
	})
}

type masterKey struct{}

// UseMaster marks the context so GetDB routes reads to the primary database.
func UseMaster(ctx context.Context) context.Context {
	return context.WithValue(ctx, masterKey{}, true)
}

// GetDB returns the transaction bound to ctx, the master database when
// UseMaster was applied, or the shared database otherwise.
func (m *DBManager) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	if ctx.Value(masterKey{}) != nil {
		return m.db.Clauses(dbresolver.Write)
	}
	return m.db
}

var _ repository.TransactionManager = (*DBManager)(nil)
