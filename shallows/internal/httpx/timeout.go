package httpx

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Timeout10s 10 second timeout handler
func Timeout10s() func(http.Handler) http.Handler {
	return TimeoutHandler(10 * time.Second)
}

// Timeout1s 1 second timeout handler
func Timeout1s() func(http.Handler) http.Handler {
	return TimeoutHandler(time.Second)
}

// Timeout2s 2 second timeout handler
func Timeout2s() func(http.Handler) http.Handler {
	return TimeoutHandler(2 * time.Second)
}

// Timeout4s 4 second timeout handler
func Timeout4s() func(http.Handler) http.Handler {
	return TimeoutHandler(4 * time.Second)
}

// UnsafeTimeout1m 1 minute timeout handler - long requests are unsafe
// and lead to system instability.
func UnsafeTimeout1m() func(http.Handler) http.Handler {
	return TimeoutHandler(time.Minute)
}

// UnsafeTimeout4m 4 minute timeout handler - long requests are unsafe
// and lead to system instability.
func UnsafeTimeout4m() func(http.Handler) http.Handler {
	return TimeoutHandler(4 * time.Minute)
}

// TimeoutHandler inserts a buffer into the http.Request context.
func TimeoutHandler(max time.Duration) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(req.Context(), max)
			defer cancel()
			original.ServeHTTP(resp, req.WithContext(ctx))
		})
	}
}

// TimeoutRollingRead specifies maximum duration between reads for a request.
func TimeoutRollingRead(max time.Duration) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithCancel(req.Context())
			defer cancel()
			req.Body = newTimeoutReader(cancel, max, req.Body)
			original.ServeHTTP(resp, req.WithContext(ctx))
		})
	}
}

// TimeoutRollingWrite specifies maximum duration between writes for a response.
func TimeoutRollingWrite(max time.Duration) func(http.Handler) http.Handler {
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			// TODO:
			// ctx, cancel := context.WithCancel(req.Context())
			// defer cancel()
			// req.Body = newTimeoutReader(cancel, max, req.Body)
			// original.ServeHTTP(resp, req.WithContext(ctx))
			original.ServeHTTP(resp, req)
		})
	}
}

func newTimeoutReader(cancel context.CancelFunc, d time.Duration, r io.ReadCloser) *timeoutreader {
	return &timeoutreader{
		cancel: cancel,
		inner:  r,
		d:      d,
		timer:  time.NewTimer(d),
	}
}

type timeoutreader struct {
	cancel context.CancelFunc
	inner  io.ReadCloser
	d      time.Duration
	timer  *time.Timer
}

func (t *timeoutreader) Read(b []byte) (n int, err error) {
	n, err = t.inner.Read(b)
	select {
	case <-t.timer.C:
		t.cancel()
		return 0, context.DeadlineExceeded
	default:
	}

	t.timer.Reset(t.d)

	return n, err
}

func (t *timeoutreader) Close() error {
	return t.inner.Close()
}
