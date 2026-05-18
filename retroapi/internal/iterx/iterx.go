package iterx

import (
	"context"
	"iter"
)

func FromChannel[T any](ch <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range ch {
			if !yield(val) {
				return // Consumer stopped iterating
			}
		}
	}
}

func From[T any](items ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, val := range items {
			if !yield(val) {
				return // Consumer stopped iterating
			}
		}
	}
}

func Chunk[T any](seq iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		chunk := make([]T, 0, n)
		for v := range seq {
			chunk = append(chunk, v)
			if len(chunk) == n {
				if !yield(chunk) {
					return
				}
				chunk = chunk[:0]
			}
		}
		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

type Seq[T any] interface {
	Each(context.Context) iter.Seq[T]
	Err() error
}
