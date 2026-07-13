package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/stretchr/testify/require"
)

func TestExtractRelease(t *testing.T) {
	tests := []struct {
		title string
		want  ddisc.ReleaseInfo
	}{
		{
			title: "Some.Movie.2024.2160p.UHD.BluRay.REMUX.HDR10.Atmos",
			want: ddisc.ReleaseInfo{
				Source:     ddisc.SourceRemux,
				Resolution: ddisc.Resolution2160p,
				HDR:        true,
				Atmos:      true,
				Remux:      true,
			},
		},
		{
			title: "Some.Movie.2024.1080p.BluRay.x264",
			want: ddisc.ReleaseInfo{
				Source:     ddisc.SourceBluRay,
				Resolution: ddisc.Resolution1080p,
			},
		},
		{
			title: "Some.Show.S01E01.720p.WEB-DL.x264",
			want: ddisc.ReleaseInfo{
				Source:     ddisc.SourceWEBDL,
				Resolution: ddisc.Resolution720p,
			},
		},
		{
			title: "Some.Show.S01E01.480p.HDTV.x264",
			want: ddisc.ReleaseInfo{
				Source:     ddisc.SourceHDTV,
				Resolution: ddisc.Resolution480p,
			},
		},
		{
			title: "totally unlabeled title",
			want: ddisc.ReleaseInfo{
				Source:     ddisc.SourceUnknown,
				Resolution: ddisc.ResolutionUnknown,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			require.Equal(t, tt.want, ddisc.ExtractRelease(tt.title))
		})
	}
}
