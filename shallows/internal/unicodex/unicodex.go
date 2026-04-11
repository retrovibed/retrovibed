package unicodex

import (
	"unicode"
)

var tightLatin = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: 0x0020, Hi: 0x007E, Stride: 1}, // Basic ASCII
		{Lo: 0x00A0, Hi: 0x024F, Stride: 1}, // Latin-1, Extended A & B
	},
}

func LowHi(table *unicode.RangeTable) (uint32, uint32) {
	if table == nil || (len(table.R16) == 0 && len(table.R32) == 0) {
		return 0, 0
	}

	var min, max uint32

	// Check 16-bit ranges
	if len(table.R16) > 0 {
		min = uint32(table.R16[0].Lo)
		max = uint32(table.R16[len(table.R16)-1].Hi)
	}

	// Check 32-bit ranges (like CJK/Emoji)
	if len(table.R32) > 0 {
		if uint32(table.R32[0].Lo) < min || min == 0 {
			min = uint32(table.R32[0].Lo)
		}
		if uint32(table.R32[len(table.R32)-1].Hi) > max {
			max = uint32(table.R32[len(table.R32)-1].Hi)
		}
	}

	return min, max
}

func ISO639_1(v string) *unicode.RangeTable {
	switch v {
	case "en", "es", "fr", "de", "it", "pt", "nl", "sv", "da", "no", "fi", "is", "vi":
		return tightLatin
	case "ru", "be", "bg", "uk", "sr", "mk", "kk", "ky", "tg":
		return unicode.Cyrillic
	case "ar", "fa", "ur", "ps", "ku":
		return unicode.Arabic
	case "he", "yi":
		return unicode.Hebrew
	case "el":
		return unicode.Greek
	case "hi", "mr", "ne", "ks":
		return unicode.Devanagari
	case "bn", "as":
		return unicode.Bengali
	case "pa":
		return unicode.Gurmukhi
	case "gu":
		return unicode.Gujarati
	case "or":
		return unicode.Oriya
	case "ta":
		return unicode.Tamil
	case "te":
		return unicode.Telugu
	case "kn":
		return unicode.Kannada
	case "ml":
		return unicode.Malayalam
	case "th":
		return unicode.Thai
	case "lo":
		return unicode.Lao
	case "ka":
		return unicode.Georgian
	case "hy":
		return unicode.Armenian
	case "zh":
		return unicode.Han // Represents CJK Ideographs
	case "ja":
		// Japanese uses Hiragana, Katakana, and Han.
		// For a simple range query, Hiragana is a good proxy.
		return unicode.Hiragana
	case "ko":
		return unicode.Hangul
	case "am":
		return unicode.Ethiopic
	case "km":
		return unicode.Khmer
	case "my":
		return unicode.Myanmar
	default:
		return nil
	}
}
