package lock

import (
	"context"
	"time"
)

// renewer runs the background toucher of a Lock: it renews the lease every
// interval until the lock is released, lost, or the survive cap is reached.
type renewer struct {
	stopCh chan struct{}
	doneCh chan struct{}
}

func (l *Lock) startTouch(interval time.Duration) {
	if interval <= 0 {
		return
	}
	r := &renewer{stopCh: make(chan struct{}), doneCh: make(chan struct{})}
	go func() {
		defer close(r.doneCh)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-r.stopCh:
				return
			case <-ticker.C:
				if !l.touchOnce() {
					return
				}
			}
		}
	}()
	l.renewer = r
}

func (l *Lock) stopTouch() {
	if l.renewer == nil {
		return
	}
	select {
	case <-l.renewer.stopCh:
	default:
		close(l.renewer.stopCh)
	}
	<-l.renewer.doneCh
	l.renewer = nil
}

// touchOnce renews the lease once. It reports whether touching should
// continue. The touch runs on a detached context: it must fire even when the
// caller's ctx is done, and opTimeout bounds each attempt.
func (l *Lock) touchOnce() bool {
	remaining := time.Until(l.surviveDeadline)
	if remaining <= 0 {
		// MaxSurviveTime reached: drop the key now instead of letting it
		// linger up to one more touch interval.
		ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
		defer cancel()
		_ = l.apply(ctx, releaseScript, l.value)
		return false
	}

	ttl := remaining
	if ttl > l.lease {
		ttl = l.lease
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	err := l.apply(ctx, refreshScript, l.value, formatMillis(ttl))
	if err == ErrLockNotHeld {
		return false
	}
	// Transient Redis errors are retried on the next tick instead of killing
	// the toucher.
	return err == nil
}
