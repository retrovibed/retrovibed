package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionPartitionAuto(t *testing.T) {
	partitions := ddisc.Partitions(16, cryptox.NewChaCha8(t.Name()))

	t.Run("assigns a partition for a resolved known media id", func(t *testing.T) {
		kid := uuid.Must(uuid.NewV7()).String()
		d := ddisc.Discovered{KnownMediaID: kid}
		ddisc.DiscoveredOptionPartitionAuto(partitions)(&d)

		require.Equal(t, partitions.Max([]byte(kid)).String(), d.Partition)
	})

	t.Run("leaves partition unchanged when known media id is nil", func(t *testing.T) {
		d := ddisc.Discovered{KnownMediaID: uuid.Nil.String(), Partition: "untouched"}
		ddisc.DiscoveredOptionPartitionAuto(partitions)(&d)

		require.Equal(t, "untouched", d.Partition)
	})

	t.Run("leaves partition unchanged when known media id is the max uuid", func(t *testing.T) {
		d := ddisc.Discovered{KnownMediaID: uuid.Max.String(), Partition: "untouched"}
		ddisc.DiscoveredOptionPartitionAuto(partitions)(&d)

		require.Equal(t, "untouched", d.Partition)
	})
}
