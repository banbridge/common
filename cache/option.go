package cache

// CacheType selects a cache backend.
type CacheType int32

const (
	// NoCache disables caching.
	NoCache CacheType = 0
	// RedisCache uses the Redis backend.
	RedisCache CacheType = 1
	// LocalCache uses the local bigcache backend.
	LocalCache CacheType = 2
	// LocalOnRedisLevel uses the two-level local-on-Redis backend.
	LocalOnRedisLevel CacheType = 3
	// CustomCache uses a caller-provided BaseCache, mainly for tests.
	CustomCache CacheType = 4
)

// CacheOption configures a CacheOptions value.
type CacheOption func(option *CacheOptions)

// CacheOptions carries backend selection for one cache operation.
type CacheOptions struct {
	CacheType       CacheType
	CustomBaseCache BaseCache
}

// MergeCacheOptions merges options into a CacheOptions value with Redis as the
// default backend.
func MergeCacheOptions(options ...CacheOption) *CacheOptions {
	opt := &CacheOptions{
		CacheType: RedisCache,
	}
	for _, o := range options {
		o(opt)
	}
	return opt
}

// UseRedis selects the Redis backend.
func UseRedis() CacheOption {
	return func(opt *CacheOptions) {
		opt.CacheType = RedisCache
	}
}

// UseLocal selects the local bigcache backend.
func UseLocal() CacheOption {
	return func(opt *CacheOptions) {
		opt.CacheType = LocalCache
	}
}

// UseLocalOnRedis selects the two-level local-on-Redis backend.
func UseLocalOnRedis() CacheOption {
	return func(opt *CacheOptions) {
		opt.CacheType = LocalOnRedisLevel
	}
}

// UseNoCache disables caching for the operation.
func UseNoCache() CacheOption {
	return func(opt *CacheOptions) {
		opt.CacheType = NoCache
	}
}

// UseCustomCache uses a caller-provided BaseCache, mainly for tests.
func UseCustomCache(bc BaseCache) CacheOption {
	return func(opt *CacheOptions) {
		opt.CacheType = CustomCache
		opt.CustomBaseCache = bc
	}
}
