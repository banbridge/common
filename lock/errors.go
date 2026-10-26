package lock

import "errors"

var (
	// ErrNotObtained is returned by Obtain when the lock could not be
	// acquired within the configured wait timeout.
	ErrNotObtained = errors.New("lock not obtained")
	// ErrLockNotHeld is returned by Release/Refresh/TTL when the caller no
	// longer holds the lock, either because it expired or was already
	// released.
	ErrLockNotHeld = errors.New("lock not held")
)
