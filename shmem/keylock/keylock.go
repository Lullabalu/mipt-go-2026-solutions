//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type KeyLock struct {
	mu     sync.Mutex
	locked map[string]bool
	wake   chan struct{}
}

func New() *KeyLock {
	return &KeyLock{
		locked: make(map[string]bool),
		wake:   make(chan struct{}),
	}
}

func (l *KeyLock) All(keys []string) bool {
	for _, k := range keys {
		if l.locked[k] {
			return false
		}
	}
	return true
}

func (l *KeyLock) Lock(keys []string) {
	for _, k := range keys {
		l.locked[k] = true
	}
}

func (l *KeyLock) unLock(keys []string) {
	for _, k := range keys {
		l.locked[k] = false
	}

	close(l.wake)
	l.wake = make(chan struct{})
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	newKeys := make([]string, len(keys))
	copy(newKeys, keys)
	sort.Strings(newKeys)

	for {
		l.mu.Lock()
		if l.All(newKeys) {
			l.Lock(newKeys)
			l.mu.Unlock()

			return false, func() {
				l.mu.Lock()
				l.unLock(newKeys)
				l.mu.Unlock()
			}
		}

		ch := l.wake
		l.mu.Unlock()

		select {
		case <-cancel:
			return true, func() {}
		case <-ch:
		}
	}
}
