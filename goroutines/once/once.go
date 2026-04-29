//go:build !solution

package once

type Once struct {
	in  chan struct{}
	out chan struct{}
}

func New() *Once {
	on := Once{in: make(chan struct{}, 1), out: make(chan struct{}, 1)}
	return &on
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
//
//	once := New()
//
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation. A new instance of
// Once is required for each function to execute.
//
// Do is intended for initialization that must be run exactly once.
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
func cl(ch chan struct{}) {
	<-ch
}

func (o *Once) Do(f func()) {
	o.in <- struct{}{}
	select {
	case o.out <- struct{}{}:
		defer cl(o.in)
		f()
		return 
	default:
		defer cl(o.in)
		return
	}
}
