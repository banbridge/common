package gormrepo

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testAccount struct {
	ID       uint   `gorm:"primaryKey,autoIncrement"`
	Address  string `gorm:"uniqueIndex"`
	Balance  float64
	Category string
}

type testGroupRow struct {
	Category string `gorm:"column:category"`
	Total    int64  `gorm:"column:total"`
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })
	require.NoError(t, db.AutoMigrate(&testAccount{}))
	return db
}
