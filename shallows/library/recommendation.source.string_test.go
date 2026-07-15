package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestRecommendationSourceString(t *testing.T) {
	t.Run("should be able to invert sources", func(t *testing.T) {
		testcase := func(s string) {
			require.Equal(t, s, library.RecommendationSourceString(md5x.String(s)))
		}
		testcase(library.RecommendationSourceRandom)
		testcase(library.RecommendationSourceDiscovered)
		testcase(library.RecommendationSourceGenerative)
	})
}
