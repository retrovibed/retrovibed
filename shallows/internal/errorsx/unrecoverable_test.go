package errorsx_test

import (
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/internal/errorsx"

	"github.com/stretchr/testify/assert"
)

func TestUnrecoverableStdLibInterop(t *testing.T) {
	var (
		local errorsx.Unrecoverable
		cause = errorsx.NewUnrecoverable(errors.New("derp"))
	)

	assert.True(t, errors.As(cause, &local))
	assert.True(t, errors.Is(cause, &local))
	assert.True(t, errors.Is(cause, local))
	assert.EqualError(t, errors.Unwrap(cause), "derp")
}

func TestUnrecoverable(t *testing.T) {
	var (
		local   errorsx.Unrecoverable
		cause   = errorsx.NewUnrecoverable(errors.New("derp"))
		wrapped = errorsx.Wrap(cause, "wrapped error")
	)

	assert.True(t, errors.As(wrapped, &local))
	assert.True(t, errors.Is(wrapped, &local))
	assert.True(t, errors.Is(wrapped, local))
}
