package lucenex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/stretchr/testify/assert"
)

func TestParsable(t *testing.T) {
	assert.True(t, lucenex.Parsable("mimetype:\"video/webm\""))
	assert.True(t, lucenex.Parsable("(mimetype:\"video/webm\" OR mimetype:\"video/ogg\")"))
	assert.True(t, lucenex.Parsable("auto_description", lucenex.WithDefaultField("auto_description")))

	assert.False(t, lucenex.Parsable(""))
	assert.False(t, lucenex.Parsable("   "))
	assert.False(t, lucenex.Parsable("\"the pitt"))
}
