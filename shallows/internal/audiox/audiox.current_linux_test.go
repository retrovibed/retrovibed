package audiox_test

import (
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/audiox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrent(t *testing.T) {
	t.Run("missing pulseaudio socket is unrecoverable", func(t *testing.T) {
		t.Setenv("PULSE_SERVER", "unix:"+t.TempDir()+"/pulse/native")

		sink, err := audiox.Current()
		require.Error(t, err)
		assert.True(t, errors.Is(err, errorsx.Unrecoverable{}))
		assert.Equal(t, audiox.Sink{}, sink)
	})
}
