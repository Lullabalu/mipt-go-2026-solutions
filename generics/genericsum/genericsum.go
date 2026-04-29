//go:build !solution

package genericsum

import (
	"slices"

	"golang.org/x/exp/constraints"
	"reflect"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[S ~[]T, T constraints.Ordered](s S) {
	slices.Sort(s)
}

func MapsEqual[Key, Value comparable](a, b map[Key]Value) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || (bv != v) {
			return false
		}
	}
	return true
}

func SliceContains[S ~[]T, T comparable](s S, t T) bool {
	for _, elem := range s {
		if elem == t {
			return true
		}
	}
	return false
}

func RunChan[T any](ch chan T, chs []<-chan T) {
	for {
		count := 0
		for _, c := range chs {
			select {
			case el, ok := <-c:
				if !ok {
					count += 1
				} else {
					ch <- el
				}
			default:
				continue
			}
		}
		if count == len(chs) {
			close(ch)
			return
		}
	}
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	result := make(chan T, 1)
	go RunChan(result, chs)
	return result
}

type Scalar interface {
	comparable
	constraints.Complex | constraints.Float | constraints.Integer
}

func conj[T Scalar](z T) T {
	v := reflect.ValueOf(any(z))

	switch v.Kind() {
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		cc := complex(real(c), -imag(c))
		return reflect.ValueOf(cc).Convert(v.Type()).Interface().(T)
	default:
		return z
	}
}

func IsHermitianMatrix[T Scalar](m [][]T) bool {
	n := len(m)
	if n == 0 {
		return true
	}
	for i := 0; i < n; i++ {
		if len(m[i]) != n {
			return false
		}
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if m[i][j] != conj(m[j][i]) {
				return false
			}
		}
	}
	return true
}
