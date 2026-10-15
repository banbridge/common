package repository

import "errors"

var (
	// ErrNotFound indicates that the target record does not exist.
	ErrNotFound = errors.New("record not found")
	// ErrConflict indicates that the operation conflicts with existing data.
	ErrConflict = errors.New("record conflict")
	// ErrDBFailed indicates that a database operation failed.
	ErrDBFailed = errors.New("database operation failed")
)

// IsNotFound reports whether err is a record-not-found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsConflict reports whether err is a unique/duplicate-key conflict.
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IgnoreNotFound returns nil when err is a record-not-found error.
func IgnoreNotFound(err error) error {
	if IsNotFound(err) {
		return nil
	}
	return err
}
