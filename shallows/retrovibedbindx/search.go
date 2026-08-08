package retrovibedbindx

import "github.com/retrovibed/retrovibed/shallows/internal/lucenex"

func Parsable(query string) bool {
	return lucenex.Parsable(query)
}
