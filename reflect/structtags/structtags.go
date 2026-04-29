//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var data sync.Map

func Unpack(req *http.Request, ptr any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	v := reflect.ValueOf(ptr).Elem()
	an, ok := data.Load(v.Type())
	var mp map[string]int
	if !ok {
		mp = make(map[string]int, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i)
			tag := fieldInfo.Tag
			name := tag.Get("http")
			if name == "" {
				name = strings.ToLower(fieldInfo.Name)
			}
			mp[name] = i
		}
		data.Store(v.Type(), mp)
	} else {
		mp = an.(map[string]int)
	}

	for name, values := range req.Form {
		i, ok := mp[name]
		if !ok {
			continue
		}
		f := v.Field(i)
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
