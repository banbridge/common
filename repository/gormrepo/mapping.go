package gormrepo

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/banbridge/common/repository"
)

// isConflictError detects backend-specific duplicate-key errors so they can be
// mapped to repository.ErrConflict at the GORM boundary.
func isConflictError(err error) bool {
	if err == nil {
		return false
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return errors.Is(err, repository.ErrConflict) ||
		errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "UNIQUE constraint") ||
		strings.Contains(err.Error(), "Duplicate entry")
}

// isNotFoundError detects both neutral and raw GORM not-found errors so they
// can be mapped to repository.ErrNotFound at the GORM boundary.
func isNotFoundError(err error) bool {
	return repository.IsNotFound(err) || errors.Is(err, gorm.ErrRecordNotFound)
}

// ignoreNotFound returns nil when err is a not-found error.
func ignoreNotFound(err error) error {
	if isNotFoundError(err) {
		return nil
	}
	return err
}
