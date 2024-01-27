package circularqueue

import (
	"errors"
	"github.com/golang-ds/queue/circularqueue"
	"github.com/google/go-cmp/cmp"
)

// CircularQueue is a fixed-sized queue that keeps returning its contents repeatedly, in order.
type CircularQueue[T any] struct {
	wrapped *queue.CircularQueue[T]
}

// New constructs a circular queue with the contents of an array or slice.
func New[T any](entries []T) CircularQueue[T] {
	wrapped := queue.New[T]()
	q := CircularQueue[T]{&wrapped}

	for _, e := range entries {
		q.wrapped.Enqueue(e)
	}

	return q
}

// SetFirst moves the indicated entry to the front of the circular queue, returning an error if not found
func (q *CircularQueue[T]) SetFirst(entry T) error {
	if q.wrapped.IsEmpty() {
		return errors.New("queue is empty")
	}

	for i := 0; i < q.wrapped.Size(); i++ {
		var first, _ = q.wrapped.First()

		if diff := cmp.Diff(entry, first); diff == "" {
			return nil // entry is found at front of queue
		}

		q.wrapped.Rotate() // rotate to the next entry and try again
	}

	return errors.New("entry not found")
}

// Next gets the next entry in the queue, wrapping around to the front if needed
func (q *CircularQueue[T]) Next() (T, error) {
	if q.wrapped.IsEmpty() {
		return *new(T), errors.New("queue is empty")
	}

	var first, ok = q.wrapped.First()
	if !ok {
		return *new(T), errors.New("entry not found")
	}

	q.wrapped.Rotate()

	return first, nil
}
