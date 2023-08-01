package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
)

func SumStrings(ss ...string) string {
	h := fnv.New64a()
	for i := range ss {
		_, _ = h.Write([]byte(ss[i]))
	}

	sum := h.Sum(nil)

	return hex.EncodeToString(sum)
}

func SumFnv64a(bs []byte) string {
	sum := fnv.New64a().Sum(bs)
	return hex.EncodeToString(sum)
}

func SumSHA256(bs []byte) string {
	sum := sha256.Sum256(bs)
	return hex.EncodeToString(sum[:])
}
