package backoffx_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/stretchr/testify/require"
)

func TestAttempt(t *testing.T) {
	const constantDelay = 5 * time.Millisecond
	s := backoffx.Constant(constantDelay)

	t.Run("timeout before first retry sleep", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context) error {
			attempt++
			return errors.ErrUnsupported
		}

		require.ErrorIs(t, backoffx.Attempt(ctx, s, doFunc), context.DeadlineExceeded)
		require.GreaterOrEqual(t, attempt, 1)
	})

	t.Run("timeout after one retry sleep", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Millisecond)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context) error {
			attempt++
			return errors.ErrUnsupported
		}
		start := time.Now()
		err := backoffx.Attempt(ctx, s, doFunc)
		elapsed := time.Since(start)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.LessOrEqual(t, 2, attempt)
		require.GreaterOrEqual(t, elapsed, constantDelay)
	})

	t.Run("timeout after three retry sleeps", func(t *testing.T) {
		const (
			delay = 50 * time.Millisecond
			es    = 3
		)
		s := backoffx.Constant(delay)
		targetTimeout := time.Duration(es)*delay + 10*time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), targetTimeout)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context) error {
			attempt++
			return errors.ErrUnsupported
		}
		start := time.Now()
		err := backoffx.Attempt(ctx, s, doFunc)
		elapsed := time.Since(start)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.GreaterOrEqual(t, es+1, attempt)
		require.GreaterOrEqual(t, elapsed, time.Duration(es)*delay)
	})
}
