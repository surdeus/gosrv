package tmplfunc

import (
	"reflect"
)

func HasField(v any, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func Sum(a, b int) int {
	return a + b
}

func Neg(a int) int {
	return -a
}

