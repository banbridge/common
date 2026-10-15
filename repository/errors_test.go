package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNotFound(t *testing.T) {
	require.True(t, IsNotFound(ErrNotFound))
	require.True(t, IsNotFound(fmt.Errorf("%w: wrapped", ErrNotFound)))
	require.False(t, IsNotFound(errors.New("other")))
	require.False(t, IsNotFound(nil))
}

func TestIsConflict(t *testing.T) {
	require.True(t, IsConflict(ErrConflict))
	require.True(t, IsConflict(fmt.Errorf("%w: wrapped", ErrConflict)))
	require.False(t, IsConflict(errors.New("connection refused")))
	require.False(t, IsConflict(nil))
}

func TestIgnoreNotFound(t *testing.T) {
	require.NoError(t, IgnoreNotFound(ErrNotFound))
	sentinel := errors.New("other")
	require.Same(t, sentinel, IgnoreNotFound(sentinel))
	require.NoError(t, IgnoreNotFound(nil))
}
