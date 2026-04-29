//go:build !solution

package dupcall

import (
	"context"
	"errors"
	"sync"
)

type Call struct {
	mu sync.Mutex

	running   bool
	cur_count int
	done      chan struct{}

	result any
	err    error
}

func Run(o *Call, ctx context.Context, cb func(context.Context) (any, error)) {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go RunFunc(ctx, cb, done, o)
	for {
		select {
		case <-done:
			o.mu.Lock()
			close(o.done)
			o.running = false
			o.mu.Unlock()

			cancel()
			return
		default:
			o.mu.Lock()
			if o.cur_count == 0 {
				cancel()
				o.mu.Unlock()
				return
			}
			o.mu.Unlock()
		}
	}

}

func RunFunc(ctx context.Context, cb func(context.Context) (any, error), done chan struct{}, o *Call) {
	o.result, o.err = cb(ctx)
	close(done)
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (any, error),
) (result any, err error) {
	o.mu.Lock()
	if o.running {
		o.cur_count += 1
		o.mu.Unlock()
		for {
			select {
			case <-ctx.Done():
				o.mu.Lock()
				o.cur_count -= 1
				o.mu.Unlock()
				return 2, errors.New("вызов был отменён")
			case <-o.done:
				o.mu.Lock()
				o.cur_count -= 1
				o.mu.Unlock()
				return o.result, o.err
			default:
				continue
			}
		}
	} else {
		o.running = true
		o.done = make(chan struct{})
		o.cur_count = 1
		o.mu.Unlock()

		cur_ctx := context.TODO()
		go Run(o, cur_ctx, cb)
		for {
			select {
			case <-ctx.Done():
				o.mu.Lock()
				o.cur_count -= 1
				o.mu.Unlock()
				return 2, errors.New("вызов был отменён")
			case <-o.done:
				o.mu.Lock()
				o.cur_count -= 1
				o.mu.Unlock()
				return o.result, o.err
			default:
				continue
			}
		}
	}
}
