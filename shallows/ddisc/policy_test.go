package ddisc_test

import (
	"math"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDefaultPolicyHardReject(t *testing.T) {
	d := ddisc.Discovered{Title: "Some.Movie.2024.HDCAM.x264"}
	require.NoError(t, ddisc.DefaultPolicy().Rank(&d))
	require.NotEmpty(t, d.PolicyRejection)
	require.EqualValues(t, math.MaxUint16, d.PolicyRank)
}

func TestDefaultPolicyQualityBracketDominatesCustomFormat(t *testing.T) {
	// lower bracket (720p WEBDL) even with every custom-format bonus stacked
	// on top must never beat a higher bracket (2160p Remux) with none.
	low := ddisc.Discovered{Title: "Some.Movie.2024.720p.WEB-DL.HDR10.Atmos.REMUX"}
	high := ddisc.Discovered{Title: "Some.Other.Movie.2024.2160p.BluRay"}

	require.NoError(t, ddisc.DefaultPolicy().Rank(&low))
	require.NoError(t, ddisc.DefaultPolicy().Rank(&high))

	require.Less(t, high.PolicyRank, low.PolicyRank, "higher quality bracket should always rank better (lower) regardless of custom-format score")
}

func TestDefaultPolicyCustomFormatBreaksTiesWithinBracket(t *testing.T) {
	plain := ddisc.Discovered{Title: "Some.Movie.2024.1080p.BluRay"}
	hdr := ddisc.Discovered{Title: "Some.Movie.2024.1080p.BluRay.HDR10"}

	require.NoError(t, ddisc.DefaultPolicy().Rank(&plain))
	require.NoError(t, ddisc.DefaultPolicy().Rank(&hdr))

	require.Less(t, hdr.PolicyRank, plain.PolicyRank, "HDR should break the tie within the same quality bracket")
}

func TestDefaultPolicyHealthMidTierDefault(t *testing.T) {
	unset := ddisc.Discovered{Title: "Some.Movie.2024.1080p.BluRay"}
	require.NoError(t, ddisc.DefaultPolicy().Rank(&unset))
	require.EqualValues(t, 25, unset.Health, "unknown health should be set to the mid-tier placeholder")

	preset := ddisc.Discovered{Title: "Some.Movie.2024.1080p.BluRay", Health: 500}
	require.NoError(t, ddisc.DefaultPolicy().Rank(&preset))
	require.EqualValues(t, 500, preset.Health, "a real pre-set health value should not be clobbered")
}
