package lock

import "time"

// LockOption configures a single Obtain call.
type LockOption func(option *lockOptions)

type lockOptions struct {
	expireTime    time.Duration
	waitTimeout   time.Duration
	retryInterval time.Duration
}

// MergeLockOptions merges options with the default lock settings: a 1 minute
// total hold cap (MaxSurviveTime), a single non-blocking acquire attempt, and
// a 50ms acquire retry interval.
func MergeLockOptions(options ...LockOption) *lockOptions {
	opt := &lockOptions{
		expireTime:    time.Minute,
		retryInterval: 50 * time.Millisecond,
	}
	for _, o := range options {
		o(opt)
	}
	return opt
}

// WithExpireTime sets MaxSurviveTime: the hard limit on the total time the
// lock can be held, counted from a successful Obtain. Defaults to 1 minute.
//
// The key itself is created with a lease of min(expireTime, 1 minute) and a
// background toucher renews it at TouchInterval (500ms, or (lease-1ms)/2 for
// leases of 1 second or less). When the cap is reached the toucher removes the
// key and stops, so a hung holder cannot keep the lock forever and later
// Release/Refresh/TTL calls report ErrLockNotHeld.
func WithExpireTime(expireTime time.Duration) LockOption {
	return func(opt *lockOptions) {
		opt.expireTime = expireTime
	}
}

// WithWaitTimeout enables blocking acquire: when the lock is taken, Obtain
// retries until the timeout elapses or the lock becomes free. Zero (default)
// means a single non-blocking attempt.
func WithWaitTimeout(timeout time.Duration) LockOption {
	return func(opt *lockOptions) {
		opt.waitTimeout = timeout
	}
}

// WithRetryInterval sets the delay between acquire attempts when blocking
// acquire is enabled via WithWaitTimeout.
func WithRetryInterval(interval time.Duration) LockOption {
	return func(opt *lockOptions) {
		opt.retryInterval = interval
	}
}
