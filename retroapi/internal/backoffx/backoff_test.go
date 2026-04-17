package backoffx

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testBackoff(t testing.TB, attempts int, s Strategy, expected ...time.Duration) {
	for i := 0; i < attempts; i++ {
		require.Equal(t, expected[i], s.Backoff(i), "attempt %d", i)
	}
}

func expectedDurationTest(t testing.TB, attempt int, s Strategy, expected time.Duration) {
	require.Equal(t, expected, s.Backoff(attempt))
}

func TestCycle(t *testing.T) {
	t.Run("more attempts than delays", func(t *testing.T) {
		testBackoff(t, 5, Cycle(1*time.Second, 2*time.Second, 3*time.Second), 1*time.Second, 2*time.Second, 3*time.Second, 1*time.Second, 2*time.Second)
	})
}

func TestExponential(t *testing.T) {
	t.Run("double each time", func(t *testing.T) {
		testBackoff(t, 5, Exponential(1*time.Second), 1*time.Second, 2*time.Second, 4*time.Second, 8*time.Second, 16*time.Second)
	})

	t.Run("should gracefully handle overflows", func(t *testing.T) {
		testBackoff(
			t,
			101,
			Exponential(1*time.Second),
			time.Second<<uint(0),
			time.Second<<uint(1),
			time.Second<<uint(2),
			time.Second<<uint(3),
			time.Second<<uint(4),
			time.Second<<uint(5),
			time.Second<<uint(6),
			time.Second<<uint(7),
			time.Second<<uint(8),
			time.Second<<uint(9),
			time.Second<<uint(10),
			time.Second<<uint(11),
			time.Second<<uint(12),
			time.Second<<uint(13),
			time.Second<<uint(14),
			time.Second<<uint(15),
			time.Second<<uint(16),
			time.Second<<uint(17),
			time.Second<<uint(18),
			time.Second<<uint(19),
			time.Second<<uint(20),
			time.Second<<uint(21),
			time.Second<<uint(22),
			time.Second<<uint(23),
			time.Second<<uint(24),
			time.Second<<uint(25),
			time.Second<<uint(26),
			time.Second<<uint(27),
			time.Second<<uint(28),
			time.Second<<uint(29),
			time.Second<<uint(30),
			time.Second<<uint(31),
			time.Second<<uint(32),
			time.Second<<uint(33),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 40
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 50
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 60
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 70
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 80
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 90
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64),
			time.Duration(math.MaxInt64), // 100
		)
	})

	t.Run("attempt 0", func(t *testing.T) {
		expectedDurationTest(t, 0, Exponential(1*time.Second), time.Duration(1*time.Second))
	})

	t.Run("attempt 1", func(t *testing.T) {
		expectedDurationTest(t, 1, Exponential(1*time.Second), time.Duration(2*time.Second))
	})

	// Here are the remaining test cases, following your provided pattern:
	t.Run("attempt 2", func(t *testing.T) {
		expectedDurationTest(t, 2, Exponential(1*time.Second), time.Duration(4*time.Second))
	})

	t.Run("attempt 3", func(t *testing.T) {
		expectedDurationTest(t, 3, Exponential(1*time.Second), time.Duration(8*time.Second))
	})

	t.Run("attempt 36", func(t *testing.T) {
		expectedDurationTest(t, 36, Exponential(1*time.Second), time.Duration(math.MaxInt64))
	})

	t.Run("attempt 37", func(t *testing.T) {
		expectedDurationTest(t, 37, Exponential(1*time.Second), time.Duration(math.MaxInt64))
	})

	t.Run("attempt 54 - overflow", func(t *testing.T) {
		expectedDurationTest(t, 54, Exponential(1*time.Second), time.Duration(math.MaxInt64))
	})

	t.Run("with scaling - attempt 0", func(t *testing.T) {
		expectedDurationTest(t, 0, Exponential(500*time.Millisecond), time.Duration(500*time.Millisecond))
	})

	t.Run("with scaling - attempt 1", func(t *testing.T) {
		expectedDurationTest(t, 1, Exponential(500*time.Millisecond), time.Duration(1*time.Second))
	})

	t.Run("with scaling - attempt 2", func(t *testing.T) {
		expectedDurationTest(t, 2, Exponential(500*time.Millisecond), time.Duration(2*time.Second))
	})

	t.Run("with scaling - attempt 3", func(t *testing.T) {
		expectedDurationTest(t, 3, Exponential(500*time.Millisecond), time.Duration(4*time.Second))
	})

	t.Run("max attempt value", func(t *testing.T) {
		expectedDurationTest(t, math.MaxInt64, Exponential(1*time.Second), time.Duration(math.MaxInt64))
	})
}

func TestConstant(t *testing.T) {
	t.Run("should remain constant", func(t *testing.T) {
		testBackoff(t, 5, Constant(1*time.Second), 1*time.Second, 1*time.Second, 1*time.Second, 1*time.Second, 1*time.Second)
	})
}

func TestAttemptV(t *testing.T) {
	const constantDelay = 5 * time.Millisecond
	s := Constant(constantDelay)

	t.Run("success on first attempt", func(t *testing.T) {
		ctx := t.Context()
		doFunc := func(_ context.Context, attempt uint) (string, error) {
			require.Equal(t, uint(0), attempt)
			return "success", nil
		}
		result, err := AttemptV(ctx, s, doFunc)
		require.NoError(t, err)
		require.Equal(t, "success", result)
	})

	t.Run("success on second attempt", func(t *testing.T) {
		ctx := t.Context()
		doFunc := func(_ context.Context, attempt uint) (string, error) {
			if attempt == 0 {
				return "", errors.ErrUnsupported
			}
			return "success", nil
		}
		result, err := AttemptV(ctx, s, doFunc)
		require.NoError(t, err)
		require.Equal(t, "success", result)
	})

	t.Run("stop attempts on first retry returns initial error", func(t *testing.T) {
		ctx := t.Context()
		initialErr := errors.New("initial error")
		doFunc := func(_ context.Context, attempt uint) (string, error) {
			if attempt == 0 {
				return "", initialErr
			}
			return "", ErrStopAttempts
		}
		result, err := AttemptV(ctx, s, doFunc)
		require.ErrorIs(t, err, initialErr)
		require.Equal(t, "", result)
	})

	t.Run("stop attempts on second retry returns previous error", func(t *testing.T) {
		ctx := t.Context()
		initialErr := errors.New("initial error")
		retryErr := errors.New("retry 1 error")
		doFunc := func(_ context.Context, attempt uint) (string, error) {
			if attempt == 0 {
				return "", initialErr
			}
			if attempt == 1 {
				return "", retryErr
			}
			return "", ErrStopAttempts
		}
		result, err := AttemptV(ctx, s, doFunc)
		require.ErrorIs(t, err, retryErr)
		require.Equal(t, "", result)
	})

	t.Run("stop attempts with wrapped error", func(t *testing.T) {
		ctx := t.Context()
		previousErr := errors.New("previous error")
		doFunc := func(_ context.Context, attempt uint) (int, error) {
			if attempt == 0 {
				return 0, errors.ErrUnsupported
			}
			if attempt == 1 {
				return 0, previousErr
			}
			return 0, errors.Join(ErrStopAttempts, errors.New("wrapped"))
		}
		result, err := AttemptV(ctx, s, doFunc)
		require.ErrorIs(t, err, previousErr)
		require.Equal(t, 0, result)
	})

	t.Run("timeout before first retry sleep", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context, uint) (string, error) {
			attempt++
			return "", errors.ErrUnsupported
		}
		time.Sleep(2 * time.Millisecond)
		_, err := AttemptV(ctx, s, doFunc)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Equal(t, 1, attempt)
	})

	t.Run("timeout after one retry sleep", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Millisecond)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context, uint) (int, error) {
			attempt++
			return 0, errors.ErrUnsupported
		}
		start := time.Now()
		_, err := AttemptV(ctx, s, doFunc)
		elapsed := time.Since(start)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.GreaterOrEqual(t, elapsed, constantDelay)
	})

	t.Run("timeout after three retry sleeps", func(t *testing.T) {
		const expectedSleeps = 3
		targetTimeout := time.Duration(expectedSleeps) * constantDelay
		ctx, cancel := context.WithTimeout(context.Background(), targetTimeout)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context, uint) (float64, error) {
			attempt++
			return 0.0, errors.ErrUnsupported
		}
		start := time.Now()
		_, err := AttemptV(ctx, s, doFunc)
		elapsed := time.Since(start)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.GreaterOrEqual(t, elapsed, targetTimeout)
	})
}

func TestAttempt(t *testing.T) {
	const constantDelay = 5 * time.Millisecond
	s := Constant(constantDelay)

	t.Run("timeout before first retry sleep", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context) error {
			attempt++
			return errors.ErrUnsupported
		}

		require.ErrorIs(t, Attempt(ctx, s, doFunc), context.DeadlineExceeded)
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
		err := Attempt(ctx, s, doFunc)
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
		s := Constant(delay)
		targetTimeout := time.Duration(es)*delay + 10*time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), targetTimeout)
		defer cancel()
		attempt := 0
		doFunc := func(context.Context) error {
			attempt++
			return errors.ErrUnsupported
		}
		start := time.Now()
		err := Attempt(ctx, s, doFunc)
		elapsed := time.Since(start)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.GreaterOrEqual(t, es+1, attempt)
		require.GreaterOrEqual(t, elapsed, time.Duration(es)*delay)
	})
}
