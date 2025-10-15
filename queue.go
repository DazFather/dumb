package main

type queue[T any] []T

func (q *queue[T]) push(items ...T) queue[T] {
	*q = append(*q, items...)
	return *q
}

func (q *queue[T]) pop() *T {
	if size := len(*q); size > 0 {
		item := []T(*q)[size-1]
		*q = []T(*q)[:size-1]
		return &item
	}
	return nil
}

func (q *queue[T]) next() *T {
	size := len(*q) - 1
	if size < 0 {
		return nil
	}

	*q = []T(*q)[:size]
	if size == 0 {
		return nil
	}

	item := []T(*q)[size-1]
	return &item
}

func (q *queue[T]) peek() *T {
	if size := len(*q); size > 0 {
		item := []T(*q)[0]
		*q = []T(*q)[1:]
		return &item
	}
	return nil
}
