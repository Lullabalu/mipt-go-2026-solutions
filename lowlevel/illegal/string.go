//go:build !solution

package illegal

import "unsafe"

func StringFromBytes(b []byte) string {
	ptr := unsafe.Pointer(&b)
	return *(*string)(ptr)
}
