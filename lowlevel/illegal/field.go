//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj any, name string, value any) {
	v := reflect.ValueOf(obj)
	v = v.Elem()
	f := v.FieldByName(name)
	ptr := unsafe.Pointer(f.UnsafeAddr())
	reflect.NewAt(f.Type(), ptr).Elem().Set(reflect.ValueOf(value))
}
