package bytespool

import (
	"bytes"
	"sync"
)

const defaultSize = 32 * 1024

type (
	Bytes       = []byte
	BytesBuffer = *bytes.Buffer
)

var gp = sync.Pool{
	New: func() any {
		buf := make(Bytes, defaultSize)
		return &buf
	},
}

// GetBytes gets a bytes buffer from the pool,
// which can specify with a size,
// default is 32k.
func GetBytes(size ...uint) Bytes {
	buf := *(gp.Get().(*Bytes))

	s := defaultSize
	if len(size) != 0 {
		s = int(size[0])
		if s == 0 {
			s = defaultSize
		}
	}
	if cap(buf) >= s {
		return buf[:s]
	}

	gp.Put(&buf)
	return make(Bytes, s)
}

// GetBuffer is similar to GetBytes,
// but it returns the bytes buffer wrapped by bytes.Buffer.
func GetBuffer(size ...uint) BytesBuffer {
	return bytes.NewBuffer(GetBytes(size...)[:0])
}

// Put puts the buffer(either Bytes or BytesBuffer) back to the pool.
func Put[T Bytes | BytesBuffer](buf T) {
	switch v := any(buf).(type) {
	case Bytes:
		gp.Put(&v)
	case BytesBuffer:
		bs := v.Bytes()
		gp.Put(&bs)
		v.Reset()
	}
}
