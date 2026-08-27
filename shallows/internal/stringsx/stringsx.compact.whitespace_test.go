package stringsx_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/stretchr/testify/assert"
)

func TestCompactWhitespace(t *testing.T) {
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("foo bar"))
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("foo   bar"))
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("foo\tbar"))
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("foo\nbar"))
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("foo \t\n  bar"))
	assert.Equal(t, "foo bar baz", stringsx.CompactWhitespace("foo  bar   baz"))
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("  foo bar  "))
	assert.Equal(t, "foo bar", stringsx.CompactWhitespace("\t\n foo bar \n\t"))
	assert.Equal(t, "", stringsx.CompactWhitespace(""))
	assert.Equal(t, "", stringsx.CompactWhitespace("   "))
	assert.Equal(t, "", stringsx.CompactWhitespace("\t\n"))
	assert.Equal(t, "foo", stringsx.CompactWhitespace("foo"))
}
