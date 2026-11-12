package tracking

import (
	"context"
	"net/url"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

type RSSOption func(*RSS)

func NewFeedRSS(id string, options ...func(*RSS)) (m RSS) {
	ts := time.Now()
	r := langx.Clone(RSS{
		ID:        id,
		NextCheck: ts,
		Digest:    uuidx.Nil().String(),
	}, langx.Compose(options...), RSSOptionDefaultEncryptionSeed)

	return r
}

func RSSOptionTestDefaults(p *RSS) {
	p.ID = langx.FirstNonZero(md5x.FormatUUID(md5x.Digest(p.URL)), uuidx.NewV4().String())
	p.EncryptionSeed = ""
	p.TTL = max(p.TTL, time.Millisecond)
	p.TTLMinimum = max(p.TTLMinimum, time.Millisecond)
	p.Digest = uuidx.Nil().String()
}

func RSSOptionNextCheck(ts time.Time) RSSOption {
	return func(r *RSS) {
		r.NextCheck = ts
	}
}

func RSSOptionDigest(s string) RSSOption {
	return func(r *RSS) {
		r.Digest = s
	}
}

// populates the encryption seed based on the url host if it hasnt already been set.
func RSSOptionDefaultEncryptionSeed(r *RSS) {
	r.EncryptionSeed = stringsx.First(
		r.EncryptionSeed, md5x.FormatUUID(
			md5x.Digest(stringsx.First(
				errorsx.Zero(url.Parse(r.URL)).Host,
				uuidx.NewV4().String())),
		),
	)
}

func RSSOptionAutoID(r *RSS) {
	r.ID = langx.FirstNonZero(r.ID, md5x.FormatUUID(md5x.Digest(r.URL)))
}

func RSSOptionURL(s string) RSSOption {
	return func(r *RSS) {
		r.URL = s
	}
}

func RSSOptionDescription(s string) RSSOption {
	return func(r *RSS) {
		r.Description = s
	}
}

func RSSOptionEncryptionSeed(s string) RSSOption {
	return func(r *RSS) {
		r.EncryptionSeed = s
	}
}

func RSSOptionAutodownload(v bool) RSSOption {
	return func(r *RSS) {
		r.Autodownload = v
	}
}

func RSSOptionAutoarchive(v bool) RSSOption {
	return func(r *RSS) {
		r.Autoarchive = v
	}
}

func RSSOptionDefaultFeeds(m RSS) RSSOption {
	// provides the values for default feeds
	return func(r *RSS) {
		*r = m
		r.ID = md5x.FormatUUID(md5x.Digest(m.URL))
		r.Digest = uuidx.Nil().String()
	}
}

func RSSQuerySearch(q string) squirrel.Sqlizer {
	return squirrelx.Noop{}
}

func RSSQueryNeedsCheck() squirrel.Sqlizer {
	return squirrel.Expr("torrents_feed_rss.next_check < NOW()")
}

func RSSSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) RSSScanner {
	return NewRSSScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func RSSSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(RSSScannerStaticColumns)...).From("torrents_feed_rss")
}
