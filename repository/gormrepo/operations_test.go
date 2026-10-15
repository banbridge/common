package gormrepo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/banbridge/common/repository"
)

func TestCRUDOperations(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	count, err := CreateModels(ctx, db,
		testAccount{Address: "abc", Balance: 100, Category: "a"},
		testAccount{Address: "def", Balance: 200, Category: "b"},
	)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	_, err = CreateModels(ctx, db, testAccount{Address: "abc", Balance: 1})
	require.Error(t, err)
	require.True(t, repository.IsConflict(err))

	count, err = UpdateModelWithConds(ctx, db, &testAccount{}, map[string]any{"balance": 150}, Equal("address", "abc"))
	require.NoError(t, err)
	require.Equal(t, 1, count)

	count, err = UpdatesModelWithValue(ctx, db, &testAccount{}, &testAccount{Balance: 250}, Equal("address", "def"))
	require.NoError(t, err)
	require.Equal(t, 1, count)

	got, err := GetModelByConds(ctx, db, &testAccount{}, Equal("address", "abc"))
	require.NoError(t, err)
	require.Equal(t, 150.0, got.Balance)

	_, err = GetModelByConds(ctx, db, &testAccount{}, Equal("address", "missing"))
	require.Error(t, err)
	require.True(t, repository.IsNotFound(err))
	require.True(t, errors.Is(err, repository.ErrNotFound))
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	rows, err := ListModelsByConds(ctx, db, &testAccount{})
	require.NoError(t, err)
	require.Len(t, rows, 2)

	ptrRows, err := ListPointerModelsByConds(ctx, db, &testAccount{})
	require.NoError(t, err)
	require.Len(t, ptrRows, 2)

	one := 1
	total, err := CountByConds(ctx, db, &testAccount{}, Offset(&one))
	require.NoError(t, err)
	require.Equal(t, int64(2), total)

	count, err = DeleteModels(ctx, db, &testAccount{}, Equal("address", "abc"))
	require.NoError(t, err)
	require.Equal(t, 1, count)

	_, err = GetModelByConds(ctx, db, &testAccount{}, Equal("address", "abc"))
	require.Error(t, err)
	require.True(t, repository.IsNotFound(err))
}

func TestBatchAndGroupOperations(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	count, err := CreateModels(ctx, db,
		testAccount{Address: "a1", Balance: 10, Category: "a"},
		testAccount{Address: "a2", Balance: 20, Category: "a"},
		testAccount{Address: "b1", Balance: 30, Category: "b"},
		testAccount{Address: "b2", Balance: 40, Category: "b"},
	)
	require.NoError(t, err)
	require.Equal(t, 4, count)

	var batched []*testAccount
	batches := 0
	processed := 0
	err = FindModelInBatches(ctx, db, &testAccount{}, &batched, 2, func(tx *gorm.DB, batch int) error {
		batches++
		processed += len(batched)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 4, processed)
	require.GreaterOrEqual(t, batches, 2)

	groups, err := GroupModelsByConds(ctx, db, &testAccount{}, testGroupRow{},
		"category, sum(balance) as total", "category")
	require.NoError(t, err)
	require.Len(t, groups, 2)

	count, err = CreateModelsWithConds(ctx, db, []*testAccount{{Address: "a1", Balance: 99}},
		clause.OnConflict{DoNothing: true})
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestFindInBatchesPropagatesCallbackError(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	_, err := CreateModels(ctx, db,
		testAccount{Address: "x1", Balance: 1},
		testAccount{Address: "x2", Balance: 2},
	)
	require.NoError(t, err)

	sentinel := errors.New("callback failed")
	var batched []*testAccount
	err = FindModelInBatches(ctx, db, &testAccount{}, &batched, 1, func(tx *gorm.DB, batch int) error {
		return sentinel
	})
	require.ErrorIs(t, err, sentinel)
}

func TestOperationsRejectNilAndEmptyInputs(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	_, err := CreateModels(ctx, nil, testAccount{Address: "nil"})
	require.Error(t, err)

	count, err := CreateModels[testAccount](ctx, db)
	require.NoError(t, err)
	require.Zero(t, count)

	_, err = CreateModelsWithConds(ctx, nil, []*testAccount{{Address: "nil"}})
	require.Error(t, err)

	count, err = CreateModelsWithConds[testAccount](ctx, db, nil)
	require.NoError(t, err)
	require.Zero(t, count)

	_, err = DeleteModels(ctx, nil, &testAccount{})
	require.Error(t, err)

	_, err = UpdateModelWithConds(ctx, db, &testAccount{}, nil)
	require.Error(t, err)

	_, err = UpdateModelWithConds(ctx, nil, &testAccount{}, map[string]any{"balance": 1})
	require.Error(t, err)

	_, err = UpdatesModelWithValue(ctx, db, &testAccount{}, nil)
	require.Error(t, err)

	_, err = GetModelByConds(ctx, nil, &testAccount{})
	require.Error(t, err)

	_, err = GetModelByConds[testAccount](ctx, db, nil)
	require.Error(t, err)

	_, err = ListModelsByConds[testAccount](ctx, db, nil)
	require.Error(t, err)

	_, err = CountByConds[testAccount](ctx, db, nil)
	require.Error(t, err)

	_, err = GroupModelsByConds[testAccount, testGroupRow](ctx, db, nil, testGroupRow{}, "category", "category")
	require.Error(t, err)

	err = FindModelInBatches(ctx, db, &testAccount{}, nil, 1, func(tx *gorm.DB, batch int) error {
		return nil
	})
	require.Error(t, err)

	err = FindModelInBatches(ctx, db, &testAccount{}, &[]*testAccount{}, 0, func(tx *gorm.DB, batch int) error {
		return nil
	})
	require.Error(t, err)

	err = FindModelInBatches(ctx, db, &testAccount{}, &[]*testAccount{}, 1, nil)
	require.Error(t, err)
}
