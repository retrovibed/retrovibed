package lucenex

import (
	"log"
	"strings"
	"unicode"

	"github.com/Masterminds/squirrel"
	"github.com/grindlemire/go-lucene"
	"github.com/grindlemire/go-lucene/pkg/lucene/expr"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/runesx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Driver interface {
	RenderParam(e *expr.Expression) (s string, params []any, err error)
}

type Option func(*config)

type config struct {
	DefaultField string
}

func WithDefaultField(s string) Option {
	return func(c *config) {
		c.DefaultField = s
	}
}

func Query(d Driver, s string, options ...Option) squirrel.Sqlizer {
	if stringsx.Blank(s) {
		return squirrelx.Noop{}
	}

	cfg := langx.Clone(config{}, options...)
	return squirrelx.SqlizerFn(func() (sql string, args []interface{}, err error) {
		ast, err := lucene.Parse(s, lucene.WithDefaultField(cfg.DefaultField))
		if err != nil {
			return "", nil, err
		}

		return d.RenderParam(ast)
	})
}

// clean a random string of text for use by lucene.
func Clean(s string) string {
	runes.Remove(runes.In(unicode.Mn))
	transformer := transform.Chain(transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)), // Filter out combining diacritical marks (Unicode category Mn - Mark, Nonspacing)
		norm.NFC,
	), runes.ReplaceIllFormed(), runesx.Replace(unicode.IsPunct, ' '), runesx.Replace(unicode.IsSymbol, ' '))
	o, _, err := transform.String(transformer, s)
	if err != nil {
		log.Println("unable to normalize lucene query", err)
		return s
	}

	o = stringsx.CompactWhitespace(o)
	o = strings.ToLower(o)
	o = strings.TrimPrefix(o, "and ")
	o = strings.ReplaceAll(o, " and ", " ")
	o = strings.TrimPrefix(o, "or ")
	o = strings.ReplaceAll(o, " or ", " ")
	o = strings.TrimPrefix(o, "to ")
	o = strings.ReplaceAll(o, " to ", " ")
	o = strings.TrimPrefix(o, "not ")
	o = strings.ReplaceAll(o, " not ", " ")

	return o
}
