// Package lock provides a Redis-based distributed lock. A lock is a single
// key acquired with SET NX PX and guarded by a random token: release and
// refresh run as Lua scripts that compare the token first, so a caller can
// never release or extend a lock taken by someone else.
package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var (
	// releaseScript deletes the key only when it still carries our token.
	releaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end`)

	// refreshScript extends the key TTL only when it still carries our token.
	refreshScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
	return 0
end`)
)

const (
	// defaultTouchInterval is the touch cadence for regular leases. Leases
	// of 2*defaultTouchInterval or less use (lease-1ms)/2 instead, so a
	// touch always lands well inside the lease.
	defaultTouchInterval = 500 * time.Millisecond
	// opTimeout bounds every single Redis operation of this package, leaving
	// room for a retry or two before the next touch is due.
	opTimeout = 200 * time.Millisecond
)

// maxLease caps a single lease (LeaseTimeout). It is a variable only so tests
// can shrink it to exercise the touch path in milliseconds.
var maxLease = time.Minute

// Locker acquires distributed locks.
type Locker interface {
	// Obtain acquires the lock for key, blocking up to the configured wait
	// timeout (WithWaitTimeout). It returns ErrNotObtained when the lock is
	// still held by someone else when the budget is exhausted.
	Obtain(ctx context.Context, key string, opts ...LockOption) (*Lock, error)
}

// RedisLocker is a Locker backed by a single Redis key. The key works on
// standalone Redis, Sentinel and Cluster alike, because one lock key always
// lives on one node.
//
// The lock runs a lease model: the key is created with a lease of
// min(expireTime, 1 minute) and a background toucher renews it until the
// total hold time reaches expireTime (MaxSurviveTime), a hard cap after which
// the lock is gone even if the holder never released it.
type RedisLocker struct {
	cli    *redis.Client
	prefix string
}

// NewRedisLocker creates a RedisLocker. prefix is prepended to every lock key
// and is expected to end with the separator used by the deployment.
func NewRedisLocker(cli *redis.Client, prefix string) Locker {
	return &RedisLocker{cli: cli, prefix: prefix}
}

func (l *RedisLocker) key(key string) string {
	return l.prefix + key
}

func (l *RedisLocker) Obtain(ctx context.Context, key string, opts ...LockOption) (*Lock, error) {
	cfg := MergeLockOptions(opts...)
	lease := cfg.expireTime
	if lease > maxLease {
		lease = maxLease
	}
	fullKey := l.key(key)
	var deadline time.Time
	if cfg.waitTimeout > 0 {
		deadline = time.Now().Add(cfg.waitTimeout)
	}

	for {
		value, err := newToken()
		if err != nil {
			return nil, err
		}
		setCtx, cancel := context.WithTimeout(ctx, opTimeout)
		ok, err := l.cli.SetNX(setCtx, fullKey, value, lease).Result()
		cancel()
		if err != nil {
			return nil, err
		}
		if ok {
			lock := &Lock{
				cli:             l.cli,
				key:             fullKey,
				value:           value,
				lease:           lease,
				surviveDeadline: time.Now().Add(cfg.expireTime),
			}
			lock.startTouch(touchInterval(lease))
			return lock, nil
		}

		if cfg.waitTimeout <= 0 || time.Now().After(deadline) {
			return nil, ErrNotObtained
		}
		select {
		case <-ctx.Done():
			return nil, ErrNotObtained
		case <-time.After(cfg.retryInterval):
		}
	}
}

// touchInterval applies the touch cadence rule: defaultTouchInterval for
// regular leases, (lease-1ms)/2 for leases of 2*defaultTouchInterval or less,
// so a touch always lands well inside the lease.
func touchInterval(lease time.Duration) time.Duration {
	if lease <= 2*defaultTouchInterval {
		return (lease - time.Millisecond) / 2
	}
	return defaultTouchInterval
}

func newToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// Lock is a distributed lock acquired via Locker. Its methods are safe to
// call after the lock was lost (expired, survive cap reached, released): they
// return ErrLockNotHeld instead of touching someone else's lock.
type Lock struct {
	cli             *redis.Client
	key             string
	value           string
	lease           time.Duration
	surviveDeadline time.Time
	renewer         *renewer
	mu              sync.Mutex
	released        bool
}

// Key returns the full Redis key of the lock, including the locker prefix.
func (l *Lock) Key() string {
	if l == nil {
		return ""
	}
	return l.key
}

// Release releases the lock. It is a no-op on a nil Lock, so callers can
// `defer lk.Release(ctx)` even when Obtain returned an error. It returns
// ErrLockNotHeld once the survive cap has passed and the toucher already
// removed the key.
func (l *Lock) Release(ctx context.Context) error {
	if l == nil {
		return nil
	}
	l.stopTouch()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.released {
		return ErrLockNotHeld
	}
	l.released = true
	return l.apply(ctx, releaseScript, l.value)
}

// Refresh extends the key TTL to ttl from now, clamped to the remaining
// survive window: no refresh can push the lock past MaxSurviveTime.
func (l *Lock) Refresh(ctx context.Context, ttl time.Duration) error {
	if l == nil {
		return ErrLockNotHeld
	}
	remaining := time.Until(l.surviveDeadline)
	if remaining <= 0 {
		return ErrLockNotHeld
	}
	if ttl > remaining {
		ttl = remaining
	}
	return l.apply(ctx, refreshScript, l.value, formatMillis(ttl))
}

// TTL reports the remaining time before the lock expires. It returns
// ErrLockNotHeld when the key is gone.
func (l *Lock) TTL(ctx context.Context) (time.Duration, error) {
	if l == nil {
		return 0, ErrLockNotHeld
	}
	ttlCtx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	d, err := l.cli.PTTL(ttlCtx, l.key).Result()
	if err != nil {
		return 0, err
	}
	if d < 0 {
		return 0, ErrLockNotHeld
	}
	return d, nil
}

// apply runs one of the token-guarded scripts and normalizes a zero result to
// ErrLockNotHeld.
func (l *Lock) apply(ctx context.Context, script *redis.Script, args ...any) error {
	opCtx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	res, err := script.Run(opCtx, l.cli, []string{l.key}, args...).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return ErrLockNotHeld
	}
	return nil
}

func formatMillis(d time.Duration) string {
	return strconv.FormatInt(d.Milliseconds(), 10)
}
