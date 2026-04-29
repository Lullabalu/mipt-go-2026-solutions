//go:build !solution

package tparallel

import "sync"

type T struct {
	parallelall chan struct{}
	parallel    bool
	cond        chan struct{}

	ready chan struct{}
	done  chan struct{}

	children sync.WaitGroup
	release  chan struct{}

	mu            sync.Mutex
	readyClosed   bool
	releaseClosed bool
}

func (t *T) closeReady() {
	t.mu.Lock()
	if !t.readyClosed {
		close(t.ready)
		t.readyClosed = true
	}
	t.mu.Unlock()
}

func (t *T) closeRelease() {
	t.mu.Lock()
	if !t.releaseClosed {
		close(t.release)
		t.releaseClosed = true
	}
	t.mu.Unlock()
}

func (t *T) Wait() {
	if t.parallel {
		if t.cond != nil {
			<-t.cond
		} else {
			<-t.parallelall
		}
	}

}

func (t *T) Parallel() {
	t.parallel = true
	t.closeReady()
	t.Wait()
}

func (t *T) Run(subtest func(t *T)) {
	new_t := T{
		parallelall: t.parallelall, parallel: false, cond: t.release, ready: make(chan struct{}),
		done: make(chan struct{}), release: make(chan struct{}),
	}

	t.children.Add(1)
	go func() {
		defer t.children.Done()
		RunTest(&new_t, subtest)
	}()

	<-new_t.ready
}

func RunTest(t *T, test func(t *T)) {
	test(t)
	if t.release != nil {
		t.closeRelease()
	}

	t.children.Wait()

	t.closeReady()
	close(t.done)
}

func Run(topTests []func(t *T)) {
	parallel := make(chan struct{})

	var tests []*T

	for _, test := range topTests {
		cur_t := T{parallelall: parallel, parallel: false, cond: nil, ready: make(chan struct{}),
			done: make(chan struct{}), release: make(chan struct{}),
		}
		go RunTest(&cur_t, test)
		<-cur_t.ready

		if cur_t.parallel {
			tests = append(tests, &cur_t)
		} else {
			<-cur_t.done
		}
	}

	close(parallel)

	for _, t := range tests {
		<-t.done
	}
}
