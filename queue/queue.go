package queue

import (
	"sync"
)

type Queue[T any] struct {
	l     sync.Mutex
	queue []*T
}

func New[T any]() *Queue[T] {
	return &(Queue[T]{
		queue: []*T{},
	})
}

func (q *Queue[T]) Enqueue(item *T) {
	q.l.Lock()
	q.queue = append(q.queue, item)
	q.l.Unlock()
}

func (q *Queue[T]) Dequeue() *T {
	if q.IsEmpty() {
		return nil
	}
	first := q.queue[0]
	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		q.queue = nil
	}
	return first
}

func (q *Queue[T]) IsEmpty() bool {
	q.l.Lock()
	empty := len(q.queue) == 0
	q.l.Unlock()
	return empty
}

func (q *Queue[T]) Peek() *T {
	if q.IsEmpty() {
		return nil
	}
	q.l.Lock()
	first := q.queue[0]
	q.l.Unlock()
	return first
}
