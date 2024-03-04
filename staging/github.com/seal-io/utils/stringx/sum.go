package stringx

import (
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
)

func SumByFNV64a(s string, ss ...string) string {
	h := fnv.New64a()

	_, _ = h.Write(ToBytes(&s))
	for i := range ss {
		_, _ = h.Write(ToBytes(&ss[i]))
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SumBytesByFNV64a(bs []byte, bss ...[]byte) string {
	h := fnv.New64a()

	_, _ = h.Write(bs)
	for i := range bss {
		_, _ = h.Write(bss[i])
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SumBySHA256(s string, ss ...string) string {
	h := sha256.New()

	_, _ = h.Write(ToBytes(&s))
	for i := range ss {
		_, _ = h.Write(ToBytes(&ss[i]))
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SumBytesBySHA256(bs []byte, bss ...[]byte) string {
	h := sha256.New()

	_, _ = h.Write(bs)
	for i := range bss {
		_, _ = h.Write(bss[i])
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SumBySHA224(s string, ss ...string) string {
	h := sha256.New224()

	_, _ = h.Write(ToBytes(&s))
	for i := range ss {
		_, _ = h.Write(ToBytes(&ss[i]))
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SumBytesBySHA224(bs []byte, bss ...[]byte) string {
	h := sha256.New224()

	_, _ = h.Write(bs)
	for i := range bss {
		_, _ = h.Write(bss[i])
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
