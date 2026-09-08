package communityapi

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestSyncPublishedBackgroundWorkerRender(t *testing.T) {
	t.Run("renders the title template from the known media", func(t *testing.T) {
		var (
			worker SyncPublishedBackgroundWorker
			pc     community.PublishedContent
			known  library.Known
			lmd    library.Metadata
		)

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults))
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, func(k *library.Known) {
			k.Title = "Blade Runner"
		}))
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))

		require.Equal(t, "Blade Runner", worker.render("{{.Known.Title}}", pc, known, lmd))
	})

	t.Run("renders the description template from every scope", func(t *testing.T) {
		var (
			worker SyncPublishedBackgroundWorker
			pc     community.PublishedContent
			known  library.Known
			lmd    library.Metadata
		)

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults, func(p *community.PublishedContent) {
			p.Title = "published title"
		}))
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, func(k *library.Known) {
			k.Overview = "a replicant hunter"
		}))
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, func(m *library.Metadata) {
			m.Description = "blade.runner.mp4"
		}))

		require.Equal(
			t,
			"published title - a replicant hunter - blade.runner.mp4",
			worker.render("{{.Published.Title}} - {{.Known.Overview}} - {{.Metadata.Description}}", pc, known, lmd),
		)
	})

	t.Run("renders static templates without any interpolation", func(t *testing.T) {
		var (
			worker SyncPublishedBackgroundWorker
			pc     community.PublishedContent
			known  library.Known
			lmd    library.Metadata
		)

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults))
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))

		require.Equal(t, "new upload", worker.render("new upload", pc, known, lmd))
	})

	t.Run("renders an empty template as blank so the caller falls back", func(t *testing.T) {
		var (
			worker SyncPublishedBackgroundWorker
			pc     community.PublishedContent
			known  library.Known
			lmd    library.Metadata
		)

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults))
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))

		require.Equal(t, "", worker.render("", pc, known, lmd))
	})

	t.Run("renders a malformed template as blank so the caller falls back", func(t *testing.T) {
		var (
			worker SyncPublishedBackgroundWorker
			pc     community.PublishedContent
			known  library.Known
			lmd    library.Metadata
		)

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults))
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))

		require.Equal(t, "", worker.render("{{.Known.Title", pc, known, lmd))
	})

	t.Run("renders a template referencing an unknown field as blank so the caller falls back", func(t *testing.T) {
		var (
			worker SyncPublishedBackgroundWorker
			pc     community.PublishedContent
			known  library.Known
			lmd    library.Metadata
		)

		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults))
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))

		require.Equal(t, "", worker.render("{{.Known.Missing}}", pc, known, lmd))
	})
}
