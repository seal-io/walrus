package strs

import (
	"reflect"
	"unsafe"
)

func FromBytes(bs *[]byte) string {
	return *(*string)(unsafe.Pointer(bs))
}

func ToBytes(s *string) (bs []byte) {
	var slice = (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	var str = (*reflect.StringHeader)(unsafe.Pointer(s))
	slice.Len = str.Len
	slice.Cap = str.Len
	slice.Data = str.Data
	return
}
