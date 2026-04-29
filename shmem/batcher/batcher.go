//go:build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync"
)

type Batcher struct {
	value *slow.Value

	mu sync.Mutex

	runb    bool
	running chan struct{}
	val     any
}

func NewBatcher(v *slow.Value) *Batcher {
	batcher := Batcher{value: v, running: make(chan struct{})}
	close(batcher.running)
	return &batcher
}

func (b *Batcher) Load() any {
	<-b.running

	b.mu.Lock()
	if !b.runb {
		b.runb = true
		b.running = make(chan struct{})
		b.mu.Unlock()

		val := b.value.Load()
		b.mu.Lock()
		b.runb = false
		b.val = val
		close(b.running)
		b.mu.Unlock()

		return val
	} else {
		b.mu.Unlock()
		<-b.running
		return b.val
	}
}
