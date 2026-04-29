//go:build !solution

package reversemap

import "reflect"

func ReverseMap(forward any) any {
	mp := reflect.ValueOf(forward)

	if mp.Kind() != reflect.Map {
		panic("это не мапа")
	}

	keyT := mp.Type().Key()
	valueT := mp.Type().Elem()

	newmptype := reflect.MapOf(valueT, keyT)
	newmp := reflect.MakeMap(newmptype)

	iter := mp.MapRange()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		newmp.SetMapIndex(value, key)
	}

	return newmp.Interface()
}
