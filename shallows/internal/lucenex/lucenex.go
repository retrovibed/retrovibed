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

func Parsable(s string, options ...Option) bool {
	if stringsx.Blank(s) {
		return false
	}

	cfg := langx.Clone(config{}, options...)
	_, err := lucene.Parse(s, lucene.WithDefaultField(cfg.DefaultField))
	return err == nil
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

		sql, args, err = d.RenderParam(ast)
		if err != nil || stringsx.Blank(sql) {
			return sql, args, err
		}

		// drivers only parenthesize the operands of an operator, never the root expression.
		// a top level OR would otherwise escape any sibling predicate it is combined with,
		// i.e. `adult = ? AND (a) OR (b)`.
		return "(" + sql + ")", args, nil
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

	words := strings.Fields(o)
	kept := words[:0]
	for _, w := range words {
		if isLuceneKeyword(w) {
			continue
		}
		kept = append(kept, w)
	}

	return strings.Join(kept, " ")
}

// isLuceneKeyword reports whether w is a lucene operator word, ignoring case.
func isLuceneKeyword(w string) bool {
	for _, k := range []string{"and", "or", "to", "not"} {
		if strings.EqualFold(w, k) {
			return true
		}
	}
	return false
}
