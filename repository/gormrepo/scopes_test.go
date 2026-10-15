package gormrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestScopeBuildersDryRun(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	buildSQL := func(scopes ...func(*gorm.DB) *gorm.DB) string {
		t.Helper()
		sqlDB := db.Session(&gorm.Session{DryRun: true})
		stmt := sqlDB.WithContext(ctx).
			Model(&testAccount{}).
			Scopes(scopes...).
			Find(&[]testAccount{}).Statement
		require.NotNil(t, stmt)
		return stmt.SQL.String()
	}

	require.NotEmpty(t, buildSQL(Clauses(clause.Locking{Strength: "UPDATE"})))
	require.NotEmpty(t, buildSQL(OrderASC("id")))
	require.NotEmpty(t, buildSQL(OrderDESC("id")))
	require.NotEmpty(t, buildSQL(Order("id")))
	require.NotEmpty(t, buildSQL(GreaterEq("balance", 1)))
	require.NotEmpty(t, buildSQL(Gte("balance", 1)))
	require.NotEmpty(t, buildSQL(EqualStrict("address", "abc")))
	require.NotEmpty(t, buildSQL(NotEqual("address", "abc")))
	require.NotEmpty(t, buildSQL(Like("address", "abc")))
	require.NotEmpty(t, buildSQL(LikeEscaped("address", "a_b%c")))
	require.NotEmpty(t, buildSQL(LikeSuffix("address", "abc")))
	require.NotEmpty(t, buildSQL(LikePatternEscaped("address", "abc")))
	require.NotEmpty(t, buildSQL(Select("id")))
	require.NotEmpty(t, buildSQL(Omit("balance")))
	require.NotEmpty(t, buildSQL(JsonColEqual("labels", "k", "v")))
}
