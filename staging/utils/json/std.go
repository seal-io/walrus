//go:build !jsoniter

package json

import (
	"encoding/json"
	"fmt"
)

var (
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

type RawMessage = json.RawMessage

// MustMarshal is similar to Marshal,
// but panics if found error.
func MustMarshal(v interface{}) []byte {
	bs, err := Marshal(v)
	if err != nil {
		panic(fmt.Errorf("error marshalling json: %w", err))
	}

	return bs
}

// MustUnmarshal is similar to Unmarshal,
// but panics if found error.
func MustUnmarshal(data []byte, v interface{}) {
	err := Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("error unmarshalling json: %w", err))
	}
}

// MustMarshalIndent is similar to MarshalIndent,
// but panics if found error.
func MustMarshalIndent(v interface{}, prefix, indent string) []byte {
	bs, err := MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(fmt.Errorf("error marshalling indent json: %w", err))
	}

	return bs
}

// ShouldMarshal is similar to Marshal,
// but never return error.
func ShouldMarshal(v interface{}) []byte {
	bs, _ := Marshal(v)
	return bs
}

// ShouldUnmarshal is similar to Unmarshal,
// but never return error.
func ShouldUnmarshal(data []byte, v interface{}) {
	_ = Unmarshal(data, v)
}

// ShouldMarshalIndent is similar to MarshalIndent,
// but never return error.
func ShouldMarshalIndent(v interface{}, prefix, indent string) []byte {
	bs, _ := MarshalIndent(v, prefix, indent)
	return bs
}
