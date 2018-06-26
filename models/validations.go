package models

import "reflect"

var lengthGreaterThanZero = func(data reflect.Value) bool {
	allowedTypes := typeList{reflect.String, reflect.Array, reflect.Slice}
	if !allowedTypes.includes(data.Kind()) {
		return false
	}

	return data.Len() != 0
}

type typeList []reflect.Kind

func (tl typeList) includes(t reflect.Kind) bool {
	for _, includedType := range tl {
		if includedType == t {
			return true
		}
	}

	return false
}
