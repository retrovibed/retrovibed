package ddisc_test

import (
	"math"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestDefaultPolicyRank(t *testing.T) {
	tests := []struct {
		name          string
		d             ddisc.Discovered
		wantRejection string
		wantRank      uint16
	}{
		{
			name:          "hdcam is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.HDCAM.x264"},
			wantRejection: `(?i)\bhdcam\b`,
			wantRank:      math.MaxUint16,
		},
		{
			name:          "cam is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.CAM.x264"},
			wantRejection: `(?i)\bcam\b`,
			wantRank:      math.MaxUint16,
		},
		{
			name:          "telesync is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.TELESYNC.x264"},
			wantRejection: `(?i)\btelesync\b`,
			wantRank:      math.MaxUint16,
		},
		{
			name:          "ts is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.TS.x264"},
			wantRejection: `(?i)\bts\b`,
			wantRank:      math.MaxUint16,
		},
		{
			name:          "telecine is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.TELECINE.x264"},
			wantRejection: `(?i)\btelecine\b`,
			wantRank:      math.MaxUint16,
		},
		{
			name:          "workprint is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.WORKPRINT.x264"},
			wantRejection: `(?i)\bworkprint\b`,
			wantRank:      math.MaxUint16,
		},
		{
			name:          "camrip is not rejected - cam term requires word boundary",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.CAMRip.x264"},
			wantRejection: "",
			wantRank:      65035,
		},
		{
			name:          "unknown bracket, no custom formats",
			d:             ddisc.Discovered{Title: "Some.Movie.2024"},
			wantRejection: "",
			wantRank:      65035,
		},
		{
			name:          "sdtv 480p is the worst known bracket",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.480p.SDTV"},
			wantRejection: "",
			wantRank:      54035,
		},
		{
			name:          "hdtv with unknown resolution",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.HDTV"},
			wantRejection: "",
			wantRank:      63035,
		},
		{
			name:          "2160p bluray, top bracket, no custom formats",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.2160p.BluRay"},
			wantRejection: "",
			wantRank:      20035,
		},
		{
			name:          "1080p bluray, no custom formats",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.1080p.BluRay"},
			wantRejection: "",
			wantRank:      30035,
		},
		{
			name:          "1080p bluray with HDR scores better within the same bracket",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.1080p.BluRay.HDR10"},
			wantRejection: "",
			wantRank:      29935,
		},
		{
			name:          "1080p web-dl screener is penalized within its bracket",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.1080p.WEB-DL.Screener"},
			wantRejection: "",
			wantRank:      31085,
		},
		{
			name:          "720p bracket with every custom-format bonus stacked never beats a bare 2160p bracket",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.720p.WEB-DL.HDR10.Atmos.REMUX"},
			wantRejection: "",
			wantRank:      38855,
		},
		{
			name:          "hdcam is hard rejected",
			d:             ddisc.Discovered{Title: "Some.Movie.2024.HDCAM.x264"},
			wantRejection: `(?i)\bhdcam\b`,
			wantRank:      math.MaxUint16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d
			require.NoError(t, ddisc.DefaultPolicy().Rank(&d))
			require.Equal(t, tt.wantRejection, d.PolicyRejection)
			require.EqualValues(t, tt.wantRank, d.PolicyRank)
		})
	}
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
