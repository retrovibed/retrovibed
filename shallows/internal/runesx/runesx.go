package runesx

import (
	"golang.org/x/text/runes"
)

// combine multiple checks into one.
func Or(tests ...func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		for _, test := range tests {
			if test(r) {
				return true
			}
		}

		return false
	}
}

func Replace(test func(rune) bool, r rune) runes.Transformer {
	return runes.Map(func(r rune) rune {
		if test(r) {
			return ' '
		}

		return r
	})
}
