package strs

import (
	"encoding/base64"
	"strings"
)

// DecodeBase64 decodes the given string, which can
// accept padded or none padded format.
func DecodeBase64(str string) (string, error) {
	// normalizes to std encoding format
	str = strings.ReplaceAll(str, "-", "+")
	str = strings.ReplaceAll(str, "_", "/")
	// normalizes to no padding format
	str = strings.TrimRight(str, "=")
	var bs, err = base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return FromBytes(&bs), nil
}

// EncodeBase64 encodes the given string,
// and then output standard format.
func EncodeBase64(src string) string {
	var enc = base64.StdEncoding
	var ret = make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(ret, ToBytes(&src))
	return FromBytes(&ret)
}
