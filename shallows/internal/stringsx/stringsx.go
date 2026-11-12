package stringsx

import (
	"regexp"
	"strings"
)

// Join a more convient string join method.
func Join(sep string, parts ...string) string {
	return strings.Join(parts, sep)
}

// Reverse returns the string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func First(values ...string) string {
	for _, s := range values {
		if len(strings.TrimSpace(s)) != 0 {
			return s
		}
	}

	return ""
}

func Blank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func Present(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

var whitespace = regexp.MustCompile(`\s+`)

func CompactWhitespace(s string) string {
	return strings.TrimSpace(whitespace.ReplaceAllString(s, " "))
}

// returns the index where the strings differ.
// or the len when they're identical
func FirstDiff(a, b string) uint {
	n := min(len(a), len(b))
	for i := range n {
		if a[i] != b[i] {
			return uint(i)
		}
	}

	return uint(n)
}
