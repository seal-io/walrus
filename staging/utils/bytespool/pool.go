package bytespool

import (
	"sync"

	"github.com/valyala/bytebufferpool"
)

const defaultBytesSliceSize = 32 * 1024

var gp = sync.Pool{
	New: func() any {
		var bs = make([]byte, defaultBytesSliceSize)
		return &bs
	},
}

func GetBuffer() *bytebufferpool.ByteBuffer {
	return bytebufferpool.Get()
}

func GetBytes(length int) []byte {
	if length <= 0 {
		length = defaultBytesSliceSize
	}
	var bsp = gp.Get().(*[]byte)
	var bs = *bsp
	if cap(bs) >= length {
		return bs[:length]
	}
	gp.Put(bsp)
	return make([]byte, length)
}

func Put(b any) {
	switch t := b.(type) {
	case []byte:
		gp.Put(&t)
	case *bytebufferpool.ByteBuffer:
		bytebufferpool.Put(t)
	}
}
