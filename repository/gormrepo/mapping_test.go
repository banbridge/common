package gormrepo

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/banbridge/common/repository"
)

func TestIsConflictError(t *testing.T) {
	require.True(t, isConflictError(&mysql.MySQLError{Number: 1062}))
	require.True(t, isConflictError(gorm.ErrDuplicatedKey))
	require.True(t, isConflictError(errors.New("Duplicate entry 'x' for key 'y'")))
	require.True(t, isConflictError(errors.New("UNIQUE constraint failed: accounts.address")))
	require.True(t, isConflictError(repository.ErrConflict))
	require.False(t, isConflictError(errors.New("connection refused")))
	require.False(t, isConflictError(nil))
}

func TestIsNotFoundError(t *testing.T) {
	require.True(t, isNotFoundError(repository.ErrNotFound))
	require.True(t, isNotFoundError(gorm.ErrRecordNotFound))
	require.False(t, isNotFoundError(errors.New("other")))
	require.False(t, isNotFoundError(nil))
}

func TestIgnoreNotFoundMapping(t *testing.T) {
	require.NoError(t, ignoreNotFound(gorm.ErrRecordNotFound))
	sentinel := errors.New("other")
	require.Same(t, sentinel, ignoreNotFound(sentinel))
	require.NoError(t, ignoreNotFound(nil))
}
