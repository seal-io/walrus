package strs

// Borrowed from github.com/thanhpk/randstr.

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
)

// list of default letters that can be used to make a random string when calling String
// function with no letters provided.
var defLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Bytes generates n random bytes.
func Bytes(n int) []byte {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return b
}

// Base64 generates a random base64 string with length of n.
func Base64(n int) string {
	return String(n, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/")
}

// Hex generates a random hex string with length of n
// e.g: 67aab2d956bd7cc621af22cfb169cba8.
func Hex(n int) string { return hex.EncodeToString(Bytes(n)) }

// String generates a random string using only letters provided in the letters parameter
// if user omit letters parameters, this function will use defLetters instead.
func String(n int, letters ...string) string {
	var (
		letterRunes []rune
		bb          bytes.Buffer
	)

	if len(letters) == 0 {
		letterRunes = defLetters
	} else {
		letterRunes = []rune(letters[0])
	}

	bb.Grow(n)

	l := uint32(len(letterRunes))
	// On each loop, generate one random rune and append to output.
	for i := 0; i < n; i++ {
		bb.WriteRune(letterRunes[binary.BigEndian.Uint32(Bytes(4))%l])
	}

	return bb.String()
}
