package gormrepo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDBManagerTransaction(t *testing.T) {
	db := newTestDB(t)
	m, err := NewManager(db)
	require.NoError(t, err)
	ctx := context.Background()

	err = m.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		_, err := CreateModels(ctx, m.GetDB(ctx), testAccount{Address: "tx", Balance: 1})
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	total, err := CountByConds(ctx, m.GetDB(ctx), &testAccount{}, Equal("address", "tx"))
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
}

func TestDBManagerTransactionRollback(t *testing.T) {
	db := newTestDB(t)
	m, err := NewManager(db)
	require.NoError(t, err)
	ctx := context.Background()

	err = m.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		_, err := CreateModels(ctx, m.GetDB(ctx), testAccount{Address: "rollback", Balance: 1})
		require.NoError(t, err)
		return errors.New("rollback")
	})
	require.Error(t, err)

	total, err := CountByConds(ctx, m.GetDB(ctx), &testAccount{}, Equal("address", "rollback"))
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
}

func TestUseTx(t *testing.T) {
	db := newTestDB(t)
	m, err := NewManager(db)
	require.NoError(t, err)
	ctx := UseTx(context.Background(), db.Begin())

	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	require.True(t, ok)
	require.Same(t, tx, m.GetDB(ctx))
}

func TestUseMasterReturnsDB(t *testing.T) {
	db := newTestDB(t)
	m, err := NewManager(db)
	require.NoError(t, err)
	ctx := UseMaster(context.Background())

	var rows []testAccount
	require.NoError(t, m.GetDB(ctx).WithContext(ctx).Find(&rows).Error)
}

func TestNewManagerRejectsNilDB(t *testing.T) {
	m, err := NewManager(nil)
	require.Error(t, err)
	require.Nil(t, m)
}
