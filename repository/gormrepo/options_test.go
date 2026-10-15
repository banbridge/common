package gormrepo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/repository"
)

func TestConvertListOptions(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	count, err := CreateModels(ctx, db,
		testAccount{Address: "a", Balance: 1},
		testAccount{Address: "b", Balance: 2},
		testAccount{Address: "c", Balance: 3},
		testAccount{Address: "d", Balance: 4},
	)
	require.NoError(t, err)
	require.Equal(t, 4, count)

	offset := 1
	limit := 2
	opts := repository.ListOptions{
		Offset:  &offset,
		Limit:   &limit,
		OrderBy: []repository.OrderBy{repository.DescBy("balance")},
	}

	var rows []testAccount
	err = db.WithContext(ctx).
		Model(&testAccount{}).
		Scopes(ConvertListOptions(opts)...).
		Find(&rows).Error
	require.NoError(t, err)
	require.Len(t, rows, 2)
	require.Equal(t, 3.0, rows[0].Balance)
	require.Equal(t, 2.0, rows[1].Balance)
}

func TestConvertListOptionsAppliesDefaultLimit(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	rowsToCreate := make([]testAccount, 0, 12)
	for i := 0; i < 12; i++ {
		rowsToCreate = append(rowsToCreate, testAccount{Address: fmt.Sprintf("addr-%d", i), Balance: float64(i)})
	}
	_, err := CreateModels(ctx, db, rowsToCreate...)
	require.NoError(t, err)

	offset := 0
	var rows []testAccount
	err = db.WithContext(ctx).Model(&testAccount{}).
		Scopes(ConvertListOptions(repository.ListOptions{Offset: &offset})...).
		Find(&rows).Error
	require.NoError(t, err)
	require.Len(t, rows, DefaultListLimit)
}

func TestConvertConditions(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	_, err := CreateModels(ctx, db,
		testAccount{Address: "abc", Balance: 100, Category: "a"},
		testAccount{Address: "def", Balance: 200, Category: "b"},
		testAccount{Address: "ghi", Balance: 300, Category: "c"},
	)
	require.NoError(t, err)

	t.Run("eq", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(ConvertConditions([]repository.Condition{{Field: "address", Operator: repository.Eq, Value: "abc"}})...).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 1)
	})

	t.Run("in and not in", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(ConvertConditions([]repository.Condition{
				{Field: "category", Operator: repository.In, Value: []string{"a", "b"}},
				{Field: "address", Operator: repository.NotIn, Value: []string{"abc"}},
			})...).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 1)
		require.Equal(t, "def", rows[0].Address)
	})

	t.Run("comparison operators", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(ConvertConditions([]repository.Condition{
				{Field: "balance", Operator: repository.Gt, Value: 100},
				{Field: "balance", Operator: repository.Lte, Value: 300},
			})...).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})

	t.Run("invalid conditions are skipped", func(t *testing.T) {
		scopes := ConvertConditions([]repository.Condition{
			{Field: "", Operator: repository.Eq, Value: "abc"},
			{Field: "address", Operator: repository.Operator("LIKE"), Value: "abc"},
		})
		require.Empty(t, scopes)
	})
}

func TestConvertConditionsStrict(t *testing.T) {
	scopes, err := ConvertConditionsStrict([]repository.Condition{
		{Field: "address", Operator: repository.Eq, Value: "abc"},
	})
	require.NoError(t, err)
	require.Len(t, scopes, 1)

	_, err = ConvertConditionsStrict([]repository.Condition{
		{Field: "", Operator: repository.Eq, Value: "abc"},
	})
	require.Error(t, err)

	_, err = ConvertConditionsStrict([]repository.Condition{
		{Field: "address", Operator: repository.Operator("LIKE"), Value: "abc"},
	})
	require.Error(t, err)
}

func TestConvertToAnySliceBranches(t *testing.T) {
	require.Equal(t, []any{}, convertToAnySlice(nil))

	ptrSlice := []string{"a", "b"}
	require.Equal(t, []any{"a", "b"}, convertToAnySlice(&ptrSlice))

	arr := [2]int{1, 2}
	require.Equal(t, []any{1, 2}, convertToAnySlice(arr))

	require.Equal(t, []any{7}, convertToAnySlice(7))
}

func TestScopeComposition(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	_, err := CreateModels(ctx, db,
		testAccount{Address: "abc", Balance: 100},
		testAccount{Address: "abd", Balance: 200},
		testAccount{Address: "xyz", Balance: 300},
	)
	require.NoError(t, err)

	t.Run("or", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(Or(Equal("address", "abc"), Equal("address", "xyz"))).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})

	t.Run("and and like", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(And(LikePrefix("address", "ab"), Gt("balance", 150))).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 1)
		require.Equal(t, "abd", rows[0].Address)
	})

	t.Run("less aliases", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(And(Lt("balance", 300), Lte("balance", 200))).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})

	t.Run("is null and is not null", func(t *testing.T) {
		var rows []testAccount
		err := db.WithContext(ctx).Model(&testAccount{}).
			Scopes(IsNotNull("address")).
			Find(&rows).Error
		require.NoError(t, err)
		require.Len(t, rows, 3)

		rows = nil
		err = db.WithContext(ctx).Model(&testAccount{}).
			Scopes(IsNull("address")).
			Find(&rows).Error
		require.NoError(t, err)
		require.Empty(t, rows)
	})
}
