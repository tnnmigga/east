package algorithms

import "container/heap"

type Heap[T any] struct {
	data []T
}

func (h *Heap[T]) Push(x T) {
	heap.Push
}
