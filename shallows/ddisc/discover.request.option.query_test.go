package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDiscoverRequestOptionQuery(t *testing.T) {
	t.Run("sets query", func(t *testing.T) {
		r := ddisc.DiscoverRequest{}
		ddisc.DiscoverRequestOptionQuery("ubuntu")(&r)
		require.Equal(t, "ubuntu", r.Query)
	})

	t.Run("overrides existing query", func(t *testing.T) {
		r := ddisc.DiscoverRequest{Query: "debian"}
		ddisc.DiscoverRequestOptionQuery("ubuntu")(&r)
		require.Equal(t, "ubuntu", r.Query)
	})
}
