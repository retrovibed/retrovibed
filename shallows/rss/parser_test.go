package rss_test

import (
	"testing"

	"github.com/james-lawrence/deeppool/internal/x/testx"
	"github.com/james-lawrence/deeppool/rss"
	"github.com/stretchr/testify/require"
)

func TestParseFixture(t *testing.T) {
	ctx, done := testx.WithDeadline(t)
	defer done()
	parsed, err := rss.Parse(ctx, testx.Read(testx.Fixture("example.1.xml")))
	require.NoError(t, err)
	require.Equal(t, len(parsed), 50)

}
