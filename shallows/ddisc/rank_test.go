package ddisc_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestSelectPicksBestCandidate(t *testing.T) {
	worst := ddisc.Discovered{ID: "worst", PolicyRank: 300}
	middle := ddisc.Discovered{ID: "middle", PolicyRank: 200}
	best := ddisc.Discovered{ID: "best", PolicyRank: 100}

	got, err := ddisc.Select(slices.Values([]ddisc.Discovered{worst, middle, best}))
	require.NoError(t, err)
	require.Equal(t, best.ID, got.ID)
}

func TestSelectSkipsRejectedCandidates(t *testing.T) {
	rejected := ddisc.Discovered{ID: "rejected", PolicyRank: 0, PolicyRejection: "cam"}
	accepted := ddisc.Discovered{ID: "accepted", PolicyRank: 100}

	got, err := ddisc.Select(slices.Values([]ddisc.Discovered{rejected, accepted}))
	require.NoError(t, err)
	require.Equal(t, accepted.ID, got.ID)
}

func TestSelectNoCandidateWhenAllRejected(t *testing.T) {
	rejected := ddisc.Discovered{ID: "rejected", PolicyRejection: "cam"}

	_, err := ddisc.Select(slices.Values([]ddisc.Discovered{rejected}))
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}

func TestSelectNoCandidateWhenEmpty(t *testing.T) {
	_, err := ddisc.Select(slices.Values([]ddisc.Discovered{}))
	require.True(t, errors.Is(err, ddisc.ErrNoCandidate))
}
