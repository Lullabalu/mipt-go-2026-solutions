//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.

type Limiter struct {
	mxcnt    int
	check    chan struct{}
	stop     bool
	interval time.Duration
	calls    []time.Time
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	lm := Limiter{mxcnt: maxCount, check: make(chan struct{}, 1), stop: false, interval: interval, calls: make([]time.Time, 0)}
	return &lm
}

func (l *Limiter) Acquire(ctx context.Context) error {
	for {
		l.check <- struct{}{}
		if l.stop {
			<-l.check
			return ErrStopped
		}

		select {
		case <-ctx.Done():
			<-l.check
			return ctx.Err()
		default:

			if len(l.calls) < l.mxcnt {
				l.calls = append(l.calls, time.Now())
				<-l.check
				return nil
			}

			now := time.Now()

			if len(l.calls) > 0 && now.Sub(l.calls[0]) >= l.interval {
				l.calls = l.calls[1:]
			}

			<-l.check
		}
	}
}

func (l *Limiter) Stop() {
	l.check <- struct{}{}
	l.stop = true
	<-l.check
}
