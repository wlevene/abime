package util

import (
	"fmt"
	"reflect"
)

func SliceInsert(ar interface{}, val interface{}, index int) interface{} {

	at := reflect.TypeOf(ar)
	vt := reflect.TypeOf(val)

	if at.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected slice, got %T", at))
	}

	if at.Elem() != vt {
		panic("first argument must be a slice of the second argument's type")
	}

	sliceVal, itemVal := reflect.ValueOf(ar), reflect.ValueOf(val)

	if index == sliceVal.Len() {
		return reflect.Append(sliceVal, itemVal).Interface()
	}

	if index > sliceVal.Len() {
		return sliceVal.Interface()
	}

	begin := sliceVal.Slice(0, index+1)
	end := sliceVal.Slice(index, sliceVal.Len())

	out := reflect.AppendSlice(begin, end)
	out.Index(index).Set(itemVal)
	return out.Interface()
}

func SliceRemove(slice interface{}, index int) interface{} {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		panic("wrong type")
	}
	resultSlice := reflect.ValueOf(slice)
	if index < 0 || index >= resultSlice.Len() {
		panic("out of bounds")
	}
	prev := resultSlice.Index(index)
	for i := index + 1; i < resultSlice.Len(); i++ {
		value := resultSlice.Index(i)
		prev.Set(value)
		prev = value
	}
	return resultSlice.Slice(0, resultSlice.Len()-1).Interface()
}
