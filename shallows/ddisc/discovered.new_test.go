package ddisc_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewDiscovered(t *testing.T) {
	t.Run("derives id and infohash from the int160", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.Equal(t, torrentx.HashUID(&id), d.ID)
		require.Equal(t, id.Bytes(), d.Infohash)
	})

	t.Run("defaults mimetype to binary", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.Equal(t, mimex.Binary, d.Mimetype)
		require.Equal(t, mimex.Application, d.Category)
		require.Equal(t, mimex.Bittorrent, d.Contentmime)
	})

	t.Run("defaults known media id and partition to the nil uuid", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.Equal(t, uuid.Nil.String(), d.KnownMediaID)
		require.Equal(t, uuid.Nil.String(), d.Partition)
	})

	t.Run("defaults locales to undetermined", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.Equal(t, language.Und.String(), d.AudioDefaultLocale)
		require.Equal(t, language.Und.String(), d.SubtitlesDefaultLocale)
	})

	t.Run("defaults created, updated, next-check, and released timestamps to neg-infinity", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.True(t, d.CreatedAt.Equal(timex.NegInf()))
		require.True(t, d.UpdatedAt.Equal(timex.NegInf()))
		require.True(t, d.NextCheckAt.Equal(timex.NegInf()))
		require.True(t, d.ReleasedAt.Equal(timex.NegInf()))
	})

	t.Run("never expires on its own: built from a confirmed infohash", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.True(t, d.TombstonedAt.Equal(timex.Inf()))
	})

	t.Run("starts in the ephemeral acquisition state", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id)

		require.EqualValues(t, ddisc.AcquisitionStateEphemeral, d.AcquisitionState)
	})

	t.Run("applies caller options on top of the defaults", func(t *testing.T) {
		id := int160.Random()
		d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionTitle("a title"))

		require.Equal(t, "a title", d.Title)
	})

	t.Run("preserves a caller-supplied released at instead of defaulting it", func(t *testing.T) {
		id := int160.Random()
		released := time.Date(2010, 5, 5, 0, 0, 0, 0, time.UTC)
		d := ddisc.NewDiscovered(&id, func(d *ddisc.Discovered) { d.ReleasedAt = released })

		require.True(t, released.Equal(d.ReleasedAt))
	})
}
