package buffer

import "context"

type BufferPool struct {
	size    int
	buffers chan *Buffer
}

func (pool *BufferPool) AcquireBuffer(ctx context.Context) (*Buffer, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b := <-pool.buffers:
		return b, nil
	}
}

func (pool *BufferPool) ReleaseBuffer(ctx context.Context, b *Buffer) error {
	b.Reset()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case pool.buffers <- b:
		return nil
	}
}

func NewBufferPool(size int) *BufferPool {
	pool := &BufferPool{
		size:    size,
		buffers: make(chan *Buffer, size),
	}

	for range size {
		pool.buffers <- New()
	}

	return pool
}
