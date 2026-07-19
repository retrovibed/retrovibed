package ddisc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareDiscoveredPolicyRankTakesPriority(t *testing.T) {
	better := Discovered{PolicyRank: 1, Health: 0, Bytes: 0}
	worse := Discovered{PolicyRank: 2, Health: 1000, Bytes: 1000}

	require.Negative(t, Compare(better, worse))
	require.Positive(t, Compare(worse, better))
}

func TestCompareDiscoveredHealthBreaksPolicyRankTie(t *testing.T) {
	better := Discovered{PolicyRank: 1, Health: 100, Bytes: 0}
	worse := Discovered{PolicyRank: 1, Health: 50, Bytes: 1000}

	require.Negative(t, Compare(better, worse))
	require.Positive(t, Compare(worse, better))
}

func TestCompareDiscoveredBytesBreaksPolicyRankAndHealthTie(t *testing.T) {
	better := Discovered{PolicyRank: 1, Health: 100, Bytes: 2000}
	worse := Discovered{PolicyRank: 1, Health: 100, Bytes: 1000}

	require.Negative(t, Compare(better, worse))
	require.Positive(t, Compare(worse, better))
}

func TestCompareDiscoveredEqual(t *testing.T) {
	a := Discovered{PolicyRank: 1, Health: 100, Bytes: 1000}
	b := Discovered{PolicyRank: 1, Health: 100, Bytes: 1000}

	require.Zero(t, Compare(a, b))
}
