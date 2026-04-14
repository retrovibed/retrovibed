package unicodex_test

import (
	"testing"
	"unicode"

	"github.com/retrovibed/retrovibed/shallows/internal/unicodex"
	"github.com/stretchr/testify/assert"
)

func TestLowHi(t *testing.T) {
	t.Run("nil table", func(t *testing.T) {
		lo, hi := unicodex.LowHi(nil)
		assert.Equal(t, uint32(0), lo)
		assert.Equal(t, uint32(0), hi)
	})

	t.Run("empty table", func(t *testing.T) {
		lo, hi := unicodex.LowHi(&unicode.RangeTable{})
		assert.Equal(t, uint32(0), lo)
		assert.Equal(t, uint32(0), hi)
	})

	t.Run("R16 only", func(t *testing.T) {
		table := &unicode.RangeTable{
			R16: []unicode.Range16{
				{Lo: 0x0041, Hi: 0x005A, Stride: 1}, // A-Z
				{Lo: 0x0061, Hi: 0x007A, Stride: 1}, // a-z
			},
		}
		lo, hi := unicodex.LowHi(table)
		assert.Equal(t, uint32(0x0041), lo)
		assert.Equal(t, uint32(0x007A), hi)
	})

	t.Run("R32 only", func(t *testing.T) {
		table := &unicode.RangeTable{
			R32: []unicode.Range32{
				{Lo: 0x1F600, Hi: 0x1F64F, Stride: 1},
			},
		}
		lo, hi := unicodex.LowHi(table)
		assert.Equal(t, uint32(0x1F600), lo)
		assert.Equal(t, uint32(0x1F64F), hi)
	})

	t.Run("R16 and R32", func(t *testing.T) {
		table := &unicode.RangeTable{
			R16: []unicode.Range16{
				{Lo: 0x0041, Hi: 0x005A, Stride: 1},
			},
			R32: []unicode.Range32{
				{Lo: 0x1F600, Hi: 0x1F64F, Stride: 1},
			},
		}
		lo, hi := unicodex.LowHi(table)
		assert.Equal(t, uint32(0x0041), lo)
		assert.Equal(t, uint32(0x1F64F), hi)
	})
}

func TestISO639_1(t *testing.T) {
	t.Run("latin scripts", func(t *testing.T) {
		for _, code := range []string{"en", "es", "fr", "de", "it", "pt", "nl", "sv", "da", "no", "fi", "is", "vi"} {
			result := unicodex.ISO639_1(code)
			assert.NotNil(t, result, "expected non-nil for %q", code)
			lo, hi := unicodex.LowHi(result)
			assert.Equal(t, uint32(0x0020), lo, "unexpected lo for %q", code)
			assert.Equal(t, uint32(0x024F), hi, "unexpected hi for %q", code)
		}
	})

	t.Run("cyrillic scripts", func(t *testing.T) {
		for _, code := range []string{"ru", "be", "bg", "uk", "sr", "mk", "kk", "ky", "tg"} {
			assert.Equal(t, unicode.Cyrillic, unicodex.ISO639_1(code), "expected Cyrillic for %q", code)
		}
	})

	t.Run("arabic scripts", func(t *testing.T) {
		for _, code := range []string{"ar", "fa", "ur", "ps", "ku"} {
			assert.Equal(t, unicode.Arabic, unicodex.ISO639_1(code), "expected Arabic for %q", code)
		}
	})

	t.Run("known single-script languages", func(t *testing.T) {
		cases := map[string]*unicode.RangeTable{
			"he": unicode.Hebrew,
			"yi": unicode.Hebrew,
			"el": unicode.Greek,
			"hi": unicode.Devanagari,
			"bn": unicode.Bengali,
			"pa": unicode.Gurmukhi,
			"gu": unicode.Gujarati,
			"ta": unicode.Tamil,
			"te": unicode.Telugu,
			"kn": unicode.Kannada,
			"ml": unicode.Malayalam,
			"th": unicode.Thai,
			"lo": unicode.Lao,
			"ka": unicode.Georgian,
			"hy": unicode.Armenian,
			"zh": unicode.Han,
			"ja": unicode.Hiragana,
			"ko": unicode.Hangul,
			"am": unicode.Ethiopic,
			"km": unicode.Khmer,
			"my": unicode.Myanmar,
		}
		for code, want := range cases {
			assert.Equal(t, want, unicodex.ISO639_1(code), "unexpected table for %q", code)
		}
	})

	t.Run("unknown language code returns nil", func(t *testing.T) {
		for _, code := range []string{"", "xx", "zz", "123"} {
			assert.Nil(t, unicodex.ISO639_1(code), "expected nil for %q", code)
		}
	})
}
