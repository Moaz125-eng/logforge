package ingest

import (
	"sync"
)

type BufferPool struct {
	pool sync.Pool
	size int
}

func NewBufferPool(size int) *BufferPool {
	return &BufferPool{
		size: size,
		pool: sync.Pool{New: func() any {
			buf := make([]byte, size)
			return &buf
		}},
	}
}

func (p *BufferPool) Get() []byte {
	raw := p.pool.Get().(*[]byte)
	return (*raw)[:0]
}

func (p *BufferPool) Put(buf []byte) {
	if cap(buf) < p.size {
		return
	}
	b := buf[:p.size]
	p.pool.Put(&b)
}

type ByteBuffer struct {
	data []byte
	pool *BufferPool
}

func (b *ByteBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *ByteBuffer) Bytes() []byte {
	return b.data
}

func (b *ByteBuffer) Release() {
	if b.pool != nil {
		b.pool.Put(b.data)
	}
	b.data = nil
}

func AcquireBuffer(pool *BufferPool) *ByteBuffer {
	return &ByteBuffer{data: pool.Get(), pool: pool}
}
