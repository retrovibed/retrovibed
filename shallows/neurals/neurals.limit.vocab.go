package neurals

import (
	"strings"

	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

// LimitVocab removes runes whose codepoint the model's modulus encoding
// cannot represent losslessly (see neurals_neural.go and
// deeppool/mediaid/v2/train.py: ord(c) % numTokens), replacing them with a
// space so word boundaries survive. Feeding these codepoints to the model
// aliases them onto unrelated tokens and has been observed to derail
// unrelated predictions later in the sequence.
func LimitVocab(s string, numTokens int64) string {
	stripped := strings.Map(func(r rune) rune {
		if int64(r) >= numTokens {
			return ' '
		}
		return r
	}, s)

	return stringsx.CompactWhitespace(stripped)
}
