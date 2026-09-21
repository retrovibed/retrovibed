package lucenex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	assert.Equal(t, "The Last of Us S01 1080p AMZN WEB DL DD 5 1 Atmos H 264 NTb", lucenex.Clean("The.Last.of.Us.S01.1080p.AMZN.WEB-DL.DD+.5.1.Atmos.H.264-NTb"))
	assert.Equal(t, "The Wrong Way Use Healing Magic S01 1080p AAC 2 0 H 265", lucenex.Clean("The.Wrong.Way.to.Use.Healing.Magic.S01.1080p.AAC.2.0.H.265"))
	assert.Equal(t, "I Am Okay With This 2020 S01 1080p DS4K NF Webrip DV HDR DDP5 1 x265 Vialle", lucenex.Clean("I Am Not Okay With This (2020) S01 (1080p DS4K NF Webrip DV HDR DDP5.1 x265) - Vialle"))
	assert.Equal(t, "foo bar", lucenex.Clean("foo.and.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("not.foo.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("and.foo.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("or.foo.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("to.foo.bar"))
	assert.Equal(t, "La Legende", lucenex.Clean("La.Légende"))
	assert.Equal(t, "foo bar", lucenex.Clean("foo.bar.and"))
	assert.Equal(t, "foo bar", lucenex.Clean("foo.bar.or"))
	assert.Equal(t, "foo bar", lucenex.Clean("foo.bar.not"))
	assert.Equal(t, "foo bar", lucenex.Clean("foo.bar.to"))
	assert.Equal(t, "Foo Bar", lucenex.Clean("Foo.AND.Bar"))
	assert.Equal(t, "Foo Bar", lucenex.Clean("NOT.Foo.Bar"))
	assert.Equal(t, "Foo Bar", lucenex.Clean("Foo.Bar.Or"))
	assert.Equal(t, "Foo Bar", lucenex.Clean("Foo.To.Bar"))
	assert.Equal(t, "How School 101 Brilliant Ideas", lucenex.Clean("How to School 101 Brilliant Ideas to"))
	assert.Equal(t, "the pitt", lucenex.Clean("\"the pitt"))
}
