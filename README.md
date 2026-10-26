# common

Shared Go foundation packages.

## Packages

- `cache` — pluggable cache facade (Redis / local bigcache / two-level) with functional options.
- `repository` — repository helpers and transaction management on top of GORM.
- `lock` — Redis-based distributed lock.
- `ctx_values` — context key helpers (log id, key-value propagation).
- `security` — crypto helpers.
- `uniq_gen` — unique id generation.

## lock

Single-key Redis lock with a lease model: `SET NX PX` to acquire, plus Lua
scripts that compare a per-lock random token before release/refresh, so a
caller can never touch a lock held by someone else. Works on standalone Redis,
Sentinel and Cluster.

Parameters:

| Parameter | Rule |
|---|---|
| LeaseTimeout | `min(expireTime, 1min)`, a single lease never exceeds 1 minute |
| TouchInterval | 500ms by default; leases of 1s or less use `(lease-1ms)/2` |
| MaxSurviveTime | equals `expireTime`, hard cap on the total hold time |
| OpTimeout | hardcoded 200ms, bounds every Redis operation |

A background toucher renews the lease until MaxSurviveTime is reached, then
removes the key and stops — a hung holder cannot keep the lock forever, and
`Release`/`Refresh`/`TTL` report `ErrLockNotHeld` afterwards. `Refresh` cannot
push the lock past the cap either.

```go
locker := lock.NewRedisLocker(redisCli, "myapp:lock:")

lk, err := locker.Obtain(ctx, "order:123",
    lock.WithExpireTime(10*time.Minute), // total hold cap, default 1 minute
    lock.WithWaitTimeout(2*time.Second))
if err != nil {
    return err // lock.ErrNotObtained when the wait budget is exhausted
}
defer lk.Release(ctx)
```

Options: `WithExpireTime` (default 1 minute), `WithWaitTimeout` (default:
single non-blocking attempt), `WithRetryInterval` (default 50ms).
