package ptr

import "reflect"

func IsPtr[T any](v T) bool {
	//if v == nil {
	//	return false
	//}

	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

func IsNil[T any](v T) bool {
	//if v == nil {
	//	return true
	//}

	if !IsPtr(v) {
		return false
	}
	return reflect.ValueOf(v).IsNil()
}
