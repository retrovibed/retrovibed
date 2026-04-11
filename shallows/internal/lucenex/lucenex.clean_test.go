package lucenex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/internal/lucenex"
	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	assert.Equal(t, "the last of us s01 1080p amzn web dl dd 5 1 atmos h 264 ntb", lucenex.Clean("The.Last.of.Us.S01.1080p.AMZN.WEB-DL.DD+.5.1.Atmos.H.264-NTb"))
	assert.Equal(t, "the wrong way use healing magic s01 1080p aac 2 0 h 265", lucenex.Clean("The.Wrong.Way.to.Use.Healing.Magic.S01.1080p.AAC.2.0.H.265"))
	assert.Equal(t, "i am okay with this 2020 s01 1080p ds4k nf webrip dv hdr ddp5 1 x265 vialle", lucenex.Clean("I Am Not Okay With This (2020) S01 (1080p DS4K NF Webrip DV HDR DDP5.1 x265) - Vialle"))
	assert.Equal(t, "foo bar", lucenex.Clean("foo.and.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("not.foo.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("and.foo.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("or.foo.bar"))
	assert.Equal(t, "foo bar", lucenex.Clean("to.foo.bar"))
	assert.Equal(t, "la legende", lucenex.Clean("La.Légende"))
}
