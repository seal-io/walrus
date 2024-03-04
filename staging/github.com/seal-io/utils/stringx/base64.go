package stringx

import (
	"encoding/base64"
	"strings"
)

// DecodeBase64 decodes the given string, which can
// accept padded or none padded format.
func DecodeBase64(str string) (string, error) {
	// Normalizes to std encoding format.
	str = strings.ReplaceAll(str, "-", "+")
	str = strings.ReplaceAll(str, "_", "/")

	// Normalizes to no padding format.
	str = strings.TrimRight(str, "=")

	bs, err := DecodeBase64Bytes(ToBytes(&str))
	if err != nil {
		return "", err
	}

	return FromBytes(&bs), nil
}

// DecodeBase64Bytes decodes the given bytes.
func DecodeBase64Bytes(src []byte) ([]byte, error) {
	enc := base64.RawStdEncoding
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)

	return dst[:n], err
}

// EncodeBase64 encodes the given string,
// and then output standard format.
func EncodeBase64(src string) string {
	bs := EncodeBase64Bytes(ToBytes(&src))
	return FromBytes(&bs)
}

// EncodeBase64Bytes is similar to DecodeBase64,
// but returns bytes.
func EncodeBase64Bytes(src []byte) []byte {
	enc := base64.RawURLEncoding
	dst := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(dst, src)

	return dst
}
