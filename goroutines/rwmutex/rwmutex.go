//go:build !solution

package rwmutex

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
type RWMutex struct {
	wrch   chan struct{}
	wach   chan struct{}
	countr int
	check  chan struct{}
}

func New() *RWMutex {
	mt := RWMutex{wrch: make(chan struct{}, 1), wach: make(chan struct{}, 1), countr: 0, check: make(chan struct{}, 1)}
	return &mt
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	for {
		rw.check <- struct{}{}
		select {
		case rw.wach <- struct{}{}:
			<-rw.wach
			select {
			case rw.wrch <- struct{}{}:
				<-rw.wrch
				rw.countr += 1
				<-rw.check
				return
			default:
				<-rw.check
				continue
			}

		default:
			<-rw.check
			continue
		}
	}
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	rw.check <- struct{}{}
	rw.countr -= 1
	<-rw.check

}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	for {
		rw.check <- struct{}{}
		select {
		case rw.wach <- struct{}{}:
			<-rw.check
			for {
				rw.check <- struct{}{}
				select {
				case rw.wrch <- struct{}{}:
					if rw.countr == 0 {
						<-rw.wach
						<-rw.check
						return
					} else {
						<-rw.wrch
						<-rw.check
						continue
					}

				default:
					<-rw.check
					continue
				}

			}

		default:
			<-rw.check
			continue
		}
	}
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	rw.check <- struct{}{}
	<-rw.wrch
	<-rw.check
}
	