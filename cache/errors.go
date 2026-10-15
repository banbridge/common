package cache

import "errors"

var (
	// ErrCacheNotFound is the normalized cache-miss error.
	ErrCacheNotFound = errors.New("cache item not found")
	// ErrNoCache indicates that caching is disabled for the operation.
	ErrNoCache = errors.New("cache is disabled")
)
