package errorsx_test

import (
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"

	"github.com/stretchr/testify/assert"
)

func TestUnrecoverable(t *testing.T) {
	t.Run("stdlib interopt", func(t *testing.T) {
		var (
			local errorsx.Unrecoverable
			cause = errorsx.NewUnrecoverable(errors.New("derp"))
		)

		assert.True(t, errors.As(cause, &local))
		assert.True(t, errors.Is(cause, &local))
		assert.True(t, errors.Is(cause, local))
		assert.EqualError(t, errors.Unwrap(cause), "derp")
	})

	t.Run("works with wrap", func(t *testing.T) {
		var (
			local   errorsx.Unrecoverable
			cause   = errorsx.NewUnrecoverable(errors.New("derp"))
			wrapped = errorsx.Wrap(cause, "wrapped error")
		)

		assert.True(t, errors.As(wrapped, &local))
		assert.True(t, errors.Is(wrapped, &local))
		assert.True(t, errors.Is(wrapped, local))
	})

}
