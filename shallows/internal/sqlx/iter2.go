package sqlx

import (
	"database/sql"
	"iter"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type scanner2[X, Y any] interface {
	Scan(x *X, y *Y) error
	Next() bool
	Close() error
	Err() error
}

type Iter2[X, Y any] interface {
	Iter() iter.Seq2[X, Y]
	Err() error
}

type scanningiter2[X, Y any] struct {
	s     scanner2[X, Y]
	cause error
}

func (t *scanningiter2[X, Y]) Iter() iter.Seq2[X, Y] {
	return func(yield func(X, Y) bool) {
		defer t.s.Close()
		for t.s.Next() {
			var (
				x X
				y Y
			)

			if err := t.s.Scan(&x, &y); err != nil {
				t.cause = errorsx.WithStack(err)
				return
			}

			if !yield(x, y) {
				return
			}
		}

		t.cause = t.s.Err()
	}
}

func (t *scanningiter2[X, Y]) Err() error {
	return t.cause
}

func Scan2[X, Y any](s scanner2[X, Y]) Iter2[X, Y] {
	return &scanningiter2[X, Y]{
		s: s,
	}
}

func Scan2One[X, Y any](s scanner2[X, Y]) (_zero0 X, _zero1 Y, _ error) {
	i := Scan2(s)

	for x, y := range i.Iter() {
		return x, y, i.Err()
	}

	return _zero0, _zero1, errorsx.Compact(i.Err(), sql.ErrNoRows)
}

func Discard2[X, Y any](s Iter2[X, Y]) (err error) {
	for range s.Iter() {
	}

	return s.Err()
}
