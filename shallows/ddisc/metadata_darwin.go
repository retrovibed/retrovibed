//go:build darwin

package ddisc

import (
	"errors"
	"io"
)

func Extract(src io.ReadSeeker) (_zero Extracted, err error) {
	return _zero, errors.ErrUnsupported
}
