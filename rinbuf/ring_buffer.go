package rinbuf

import "fmt"

type RingBuffer[T any] struct {
	buffer []T
	next   int
}

func New[T any](size int) RingBuffer[T] {
	return RingBuffer[T]{
		buffer: make([]T, size),
		next:   0,
	}
}

func (buf *RingBuffer[T]) Push(val T) (removed T) {
	removed = buf.buffer[buf.next]
	buf.buffer[buf.next] = val
	buf.next = (buf.next + 1) % len(buf.buffer)
	return removed
}

func (buf *RingBuffer[T]) Len() int {
	return len(buf.buffer)
}

func (buf *RingBuffer[T]) Peek(idx int) T {
	return buf.buffer[(buf.next+idx)%len(buf.buffer)]
}

func (buf *RingBuffer[T]) Buffer() []T {
	return buf.buffer
}

func (buf *RingBuffer[T]) Resize(size int) error {
	if size > len(buf.buffer) {
		buf.Grow(size - len(buf.buffer))
	} else if size < len(buf.buffer) {
		_, err := buf.Shrink(len(buf.buffer) - size)
		return err
	}

	return nil
}

func (buf *RingBuffer[T]) Grow(amount int) (added []T) {
	added = make([]T, amount)
	tail := append(added, buf.buffer[buf.next:]...)
	buf.buffer = append(buf.buffer[:buf.next], tail...)
	return added
}

func (buf *RingBuffer[T]) Shrink(amount int) (removed []T, err error) {
	if len(buf.buffer)-amount < 0 {
		return nil, fmt.Errorf("cannot shrink buffer by %d, would be negative", amount)
	}

	removed = make([]T, amount)
	if buf.next >= len(buf.buffer) {
		buf.next = 0
	}
	for i := 0; i < amount; i++ {
		removed[i] = buf.removeNext()
		if buf.next >= len(buf.buffer) {
			buf.next = 0
		}
	}

	return removed, nil
}

func (buf *RingBuffer[T]) removeNext() (removed T) {
	removed = buf.buffer[buf.next]
	buf.buffer = append(buf.buffer[:buf.next], buf.buffer[buf.next+1:]...)
	return removed
}
