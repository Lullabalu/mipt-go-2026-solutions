//go:build !solution

package jsonlist

import (
	"encoding/json"
	"io"
	"reflect"
)

func Marshal(w io.Writer, slice any) error {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}

	space := []byte{' '}
	for i := range v.Len() {
		el := v.Index(i)

		bytes, err := json.Marshal(el.Interface())

		if err != nil {
			return err
		}
		_, err = w.Write(bytes)
		if err != nil {
			return err
		}
		if i != v.Len()-1 {
			_, err = w.Write(space)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func Unmarshal(r io.Reader, slice any) error {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}

	dec := json.NewDecoder(r)

	sv := v.Elem()
	elemT := sv.Type().Elem()
	for {
		elemPtr := reflect.New(elemT)

		err := dec.Decode(elemPtr.Interface())
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		sv.Set(reflect.Append(sv, elemPtr.Elem()))
	}

	return nil
}
