//go:build jsoniter

package json

import (
	stdjson "encoding/json"
	"fmt"
	"strconv"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var (
	_json         = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = _json.Marshal
	Unmarshal     = _json.Unmarshal
	MarshalIndent = _json.MarshalIndent
	NewDecoder    = _json.NewDecoder
	NewEncoder    = _json.NewEncoder
)

func init() {
	// borrowed from https://github.com/json-iterator/go/issues/145#issuecomment-323483602
	decodeNumberAsInt64IfPossible := func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.NumberValue:
			var number stdjson.Number

			iter.ReadVal(&number)
			i, err := strconv.ParseInt(string(number), 10, 64)

			if err == nil {
				*(*any)(ptr) = i
				return
			}

			f, err := strconv.ParseFloat(string(number), 64)
			if err == nil {
				*(*any)(ptr) = f
				return
			}
		default:
			*(*any)(ptr) = iter.Read()
		}
	}
	jsoniter.RegisterTypeDecoderFunc("interface {}", decodeNumberAsInt64IfPossible)
}

// MustMarshal is similar to Marshal,
// but panics if found error.
func MustMarshal(v any) []byte {
	bs, err := Marshal(v)
	if err != nil {
		panic(fmt.Errorf("marshal json: %w", err))
	}

	return bs
}

// MustUnmarshal is similar to Unmarshal,
// but panics if found error.
func MustUnmarshal(data []byte, v any) {
	err := Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("unmarshal json: %w", err))
	}
}

// MustMarshalIndent is similar to MarshalIndent,
// but panics if found error.
func MustMarshalIndent(v any, prefix, indent string) []byte {
	bs, err := MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(fmt.Errorf("marshal indent json: %w", err))
	}

	return bs
}

// ShouldMarshal is similar to Marshal,
// but never return error.
func ShouldMarshal(v any) []byte {
	bs, _ := Marshal(v)
	return bs
}

// ShouldUnmarshal is similar to Unmarshal,
// but never return error.
func ShouldUnmarshal(data []byte, v any) {
	_ = Unmarshal(data, v)
}

// ShouldMarshalIndent is similar to MarshalIndent,
// but never return error.
func ShouldMarshalIndent(v any, prefix, indent string) []byte {
	bs, _ := MarshalIndent(v, prefix, indent)
	return bs
}
