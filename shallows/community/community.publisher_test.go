package community_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityPublisher(t *testing.T) {
	t.Run("insert with defaults properly updates fields", func(t *testing.T) {
		var (
			b, a community.CommunityPublisher
		)

		q := sqltestx.Metadatabase(t)
		require.NoError(t, testx.Fake(&b, community.CommunityPublisherOptionTestDefaults))
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(t.Context(), sqlx.Debug(q), b).Scan(&b))
		assert.EqualValues(t, uuidx.WithSuffix(1), b.CommunityID)
		assert.EqualValues(t, uuidx.WithSuffix(2), b.PublisherID)

		b.TemplateTitle = "Derp 0"
		b.TemplateDescription = "Derp 1"
		require.NoError(t, community.CommunityPublisherInsertWithDefaults(t.Context(), sqlx.Debug(q), b).Scan(&a))

		assert.EqualValues(t, "Derp 0", a.TemplateTitle)
		assert.EqualValues(t, "Derp 1", a.TemplateDescription)
	})
}
