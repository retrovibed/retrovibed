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

type Seq[T any] interface {
	Each(context.Context) iter.Seq[T]
	Err() error
}
