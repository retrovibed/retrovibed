package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredOptionDetectCorrupted(t *testing.T) {
	t.Run("valid utf8 title is left untouched", func(t *testing.T) {
		d := ddisc.Discovered{Title: "a perfectly fine title", SyncUID: "keep-me"}
		ddisc.DiscoveredOptionDetectCorrupted(&d)

		require.Equal(t, "a perfectly fine title", d.Title)
		require.Equal(t, "keep-me", d.SyncUID)
	})

	t.Run("invalid utf8 title is sanitized and excluded from sync", func(t *testing.T) {
		// strings.ToValidUTF8 collapses the contiguous \xff\xfe run into a
		// single replacement glyph, not one per invalid byte.
		d := ddisc.Discovered{Title: "broken\xff\xfename", SyncUID: "some-sync-id"}
		ddisc.DiscoveredOptionDetectCorrupted(&d)

		require.Equal(t, "broken�name", d.Title)
		require.Equal(t, uuid.Nil.String(), d.SyncUID)
	})
}
