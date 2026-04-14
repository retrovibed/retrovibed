package sqlx

import (
	"database/sql"
	"iter"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type scanner[T any] interface {
	Scan(i *T) error
	Next() bool
	Close() error
	Err() error
}

type Iter[T any] interface {
	Iter() iter.Seq[T]
	Err() error
}

type scanningiter[T any] struct {
	s     scanner[T]
	cause error
}

func (t *scanningiter[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		defer t.s.Close()
		for t.s.Next() {
			var (
				p T
			)

			if err := t.s.Scan(&p); err != nil {
				t.cause = errorsx.WithStack(err)
				return
			}

			if !yield(p) {
				return
			}
		}

		t.cause = t.s.Err()
	}
}

func (t *scanningiter[T]) Err() error {
	return t.cause
}

func Scan[T any](s scanner[T]) Iter[T] {
	return &scanningiter[T]{
		s: s,
	}
}

func ScanOne[T any](s scanner[T]) (_zero T, _ error) {
	i := Scan(s)

	for v := range i.Iter() {
		return v, i.Err()
	}

	return _zero, errorsx.Compact(i.Err(), sql.ErrNoRows)
}

// ScanInto a slice, automatically closes the scanner once done.
func ScanInto[T any](s scanner[T], dst *[]T) (err error) {
	i := Scan(s)
	for v := range i.Iter() {
		*dst = append(*dst, v)
	}

	return i.Err()
}

func Discard[T any](s Iter[T]) (err error) {
	for range s.Iter() {
	}

	return s.Err()
}
