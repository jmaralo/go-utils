// Package rinbuf provides a generic ring buffer implementation.
package rinbuf

import "fmt"

type RingBuffer[T any] struct {
	buffer []T
	next   int
}

// Create a new RingBuffer with the given size.
func New[T any](size int) RingBuffer[T] {
	return RingBuffer[T]{
		buffer: make([]T, size),
		next:   0,
	}
}

// Add a new value to the ring buffer. Returns the value that was removed.
func (buf *RingBuffer[T]) Push(val T) (removed T) {
	removed = buf.buffer[buf.next]
	buf.buffer[buf.next] = val
	buf.next = (buf.next + 1) % len(buf.buffer)
	return removed
}

// Returns the length of the ring buffer.
func (buf *RingBuffer[T]) Len() int {
	return len(buf.buffer)
}

// Returns the value that is n elements ahead of the current pointer.
func (buf *RingBuffer[T]) Peek(n int) T {
	return buf.buffer[(buf.next+n)%len(buf.buffer)]
}

// Returns the underlying buffer.
func (buf *RingBuffer[T]) Buffer() []T {
	return buf.buffer
}

// Resize the buffer to the given size, growing or shrinking as necessary.
// Returns an error if the new size is less than or equal to zero.
func (buf *RingBuffer[T]) Resize(size int) error {
	if size > len(buf.buffer) {
		buf.Grow(size - len(buf.buffer))
	} else if size < len(buf.buffer) {
		_, err := buf.Shrink(len(buf.buffer) - size)
		return err
	}

	return nil
}

// Increase the buffer size by the given amount. Returning the elements that
// were added.
func (buf *RingBuffer[T]) Grow(amount int) (added []T) {
	added = make([]T, amount)
	tail := append(added, buf.buffer[buf.next:]...)
	buf.buffer = append(buf.buffer[:buf.next], tail...)
	return added
}

// Decrease the buffer size by the given amount. Returning the elements that
// were removed.
func (buf *RingBuffer[T]) Shrink(amount int) (removed []T, err error) {
	if len(buf.buffer)-amount <= 0 {
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
