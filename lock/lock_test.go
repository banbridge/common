package lock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func setupLocker(t *testing.T) (Locker, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	// Expire keys by wall clock, so the real-time sleeps below observe TTLs.
	mr.ClockTTL(true)
	cli := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = cli.Close() })
	return NewRedisLocker(cli, "test:lock:"), mr
}

func TestObtainAndRelease(t *testing.T) {
	ctx := context.Background()
	locker, mr := setupLocker(t)

	lk, err := locker.Obtain(ctx, "order:1")
	require.NoError(t, err)
	require.Equal(t, "test:lock:order:1", lk.Key())

	_, err = locker.Obtain(ctx, "order:1")
	require.ErrorIs(t, err, ErrNotObtained)

	require.NoError(t, lk.Release(ctx))
	require.False(t, mr.Exists("test:lock:order:1"), "lock key should be deleted on release")

	lk2, err := locker.Obtain(ctx, "order:1")
	require.NoError(t, err)
	require.NoError(t, lk2.Release(ctx))
}

func TestReleaseTwice(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	lk, err := locker.Obtain(ctx, "k")
	require.NoError(t, err)
	require.NoError(t, lk.Release(ctx))
	require.ErrorIs(t, lk.Release(ctx), ErrLockNotHeld)

	var nilLock *Lock
	require.NoError(t, nilLock.Release(ctx))
}

func TestLeaseCappedAtOneMinute(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	// expireTime 2min > maxLease: the key TTL must be capped at 1 minute
	// even though the total survive window is 2 minutes.
	lk, err := locker.Obtain(ctx, "k", WithExpireTime(2*time.Minute))
	require.NoError(t, err)
	defer lk.Release(ctx)

	remaining, err := lk.TTL(ctx)
	require.NoError(t, err)
	require.Greater(t, remaining, 50*time.Second)
	require.LessOrEqual(t, remaining, time.Minute)
}

func TestRefreshClampedToSurviveCap(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	lk, err := locker.Obtain(ctx, "k", WithExpireTime(30*time.Second))
	require.NoError(t, err)

	// No refresh may push the lock past MaxSurviveTime (30s here).
	require.NoError(t, lk.Refresh(ctx, 10*time.Minute))
	remaining, err := lk.TTL(ctx)
	require.NoError(t, err)
	require.Greater(t, remaining, time.Duration(0))
	require.LessOrEqual(t, remaining, 30*time.Second)

	require.NoError(t, lk.Release(ctx))
	require.ErrorIs(t, lk.Refresh(ctx, time.Minute), ErrLockNotHeld)
}

func TestExpiredLockIsNotHeld(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	lk, err := locker.Obtain(ctx, "k", WithExpireTime(100*time.Millisecond))
	require.NoError(t, err)

	// Sleep well past the TTL so a slow scheduler cannot make this flaky.
	time.Sleep(400 * time.Millisecond)
	_, err = lk.TTL(ctx)
	require.ErrorIs(t, err, ErrLockNotHeld)
	require.ErrorIs(t, lk.Release(ctx), ErrLockNotHeld)

	// The expired lock can be acquired again.
	lk2, err := locker.Obtain(ctx, "k")
	require.NoError(t, err)
	require.NoError(t, lk2.Release(ctx))
}

func TestTouchKeepsLockAliveUntilCap(t *testing.T) {
	ctx := context.Background()
	locker, mr := setupLocker(t)

	// Shrink the lease cap so the touch path is observable in milliseconds:
	// lease becomes 200ms, touch interval (200ms-1ms)/2 ~ 99ms.
	oldMaxLease := maxLease
	maxLease = 200 * time.Millisecond
	t.Cleanup(func() { maxLease = oldMaxLease })

	lk, err := locker.Obtain(ctx, "k", WithExpireTime(1*time.Second))
	require.NoError(t, err)

	// At 700ms an untouched 200ms lease would be long gone; the toucher must
	// have kept the key alive.
	time.Sleep(700 * time.Millisecond)
	require.True(t, mr.Exists("test:lock:k"))
	remaining, err := lk.TTL(ctx)
	require.NoError(t, err)
	require.Greater(t, remaining, time.Duration(0))
	require.LessOrEqual(t, remaining, 200*time.Millisecond)

	// MaxSurviveTime (1s) is a hard cap: shortly after it the key is gone
	// even though nobody released the lock.
	time.Sleep(600 * time.Millisecond)
	require.False(t, mr.Exists("test:lock:k"))
	require.ErrorIs(t, lk.Release(ctx), ErrLockNotHeld)
}

func TestTouchStopsWhenLockLost(t *testing.T) {
	ctx := context.Background()
	locker, mr := setupLocker(t)

	lk, err := locker.Obtain(ctx, "k", WithExpireTime(time.Minute))
	require.NoError(t, err)

	// The key disappears behind the holder's back (ops delete, Redis failover
	// with data loss, ...) and someone else acquires a fresh lock.
	mr.Del("test:lock:k")
	other, err := locker.Obtain(ctx, "k", WithExpireTime(time.Minute))
	require.NoError(t, err)

	// The toucher of the original holder must stop instead of renewing, and
	// the original holder must never touch the new holder's lock. Sleep past
	// one touch interval (500ms) so at least one touch attempt happens.
	time.Sleep(600 * time.Millisecond)
	require.True(t, mr.Exists("test:lock:k"))
	require.ErrorIs(t, lk.Release(ctx), ErrLockNotHeld)
	require.NoError(t, other.Release(ctx))
}

func TestTouchIntervalRule(t *testing.T) {
	require.Equal(t, 500*time.Millisecond, touchInterval(2*time.Minute))
	require.Equal(t, 500*time.Millisecond, touchInterval(1001*time.Millisecond))
	require.Equal(t, (1000*time.Millisecond-time.Millisecond)/2, touchInterval(1000*time.Millisecond))
	require.Equal(t, (300*time.Millisecond-time.Millisecond)/2, touchInterval(300*time.Millisecond))
	require.Equal(t, (100*time.Millisecond-time.Millisecond)/2, touchInterval(100*time.Millisecond))
}

func TestBlockingAcquire(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	holder, err := locker.Obtain(ctx, "k")
	require.NoError(t, err)

	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = holder.Release(ctx)
	}()

	start := time.Now()
	lk, err := locker.Obtain(ctx, "k",
		WithWaitTimeout(2*time.Second),
		WithRetryInterval(20*time.Millisecond))
	require.NoError(t, err)
	require.GreaterOrEqual(t, time.Since(start), 100*time.Millisecond)
	require.NoError(t, lk.Release(ctx))
}

func TestWaitTimeoutExhausted(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	holder, err := locker.Obtain(ctx, "k")
	require.NoError(t, err)
	defer holder.Release(ctx)

	start := time.Now()
	_, err = locker.Obtain(ctx, "k",
		WithWaitTimeout(150*time.Millisecond),
		WithRetryInterval(20*time.Millisecond))
	require.ErrorIs(t, err, ErrNotObtained)
	require.GreaterOrEqual(t, time.Since(start), 150*time.Millisecond)
}

func TestContextCancelStopsWaiting(t *testing.T) {
	ctx := context.Background()
	locker, _ := setupLocker(t)

	holder, err := locker.Obtain(ctx, "k")
	require.NoError(t, err)
	defer holder.Release(ctx)

	waitCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	_, err = locker.Obtain(waitCtx, "k", WithWaitTimeout(10*time.Second))
	// If the ctx is already done before the first attempt, Obtain surfaces
	// the raw context error; otherwise waiting stops with ErrNotObtained.
	require.True(t,
		errors.Is(err, ErrNotObtained) || errors.Is(err, context.DeadlineExceeded),
		"unexpected error: %v", err)
}

func TestPrefixIsolation(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	cli := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = cli.Close() })

	app1 := NewRedisLocker(cli, "app1:lock:")
	app2 := NewRedisLocker(cli, "app2:lock:")

	lk1, err := app1.Obtain(ctx, "k")
	require.NoError(t, err)

	lk2, err := app2.Obtain(ctx, "k")
	require.NoError(t, err)

	require.NoError(t, lk1.Release(ctx))
	require.NoError(t, lk2.Release(ctx))
}
