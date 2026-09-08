package backups_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuetestx"
	"github.com/stretchr/testify/require"
)

func TestEnqueue(t *testing.T) {
	t.Run("queues a request for the device", func(t *testing.T) {
		wq := pqueuetestx.NewQueue()

		require.NoError(t, backups.Enqueue(t.Context(), wq, "device-1"))
		require.Equal(t, 1, wq.Len())
		require.JSONEq(t, `{"device":"device-1"}`, string(wq.Snapshot()[0]))
	})

	t.Run("queues every request", func(t *testing.T) {
		wq := pqueuetestx.NewQueue()

		require.NoError(t, backups.Enqueue(t.Context(), wq, "device-1"))
		require.NoError(t, backups.Enqueue(t.Context(), wq, "device-1"))
		require.Equal(t, 2, wq.Len())
	})
}
