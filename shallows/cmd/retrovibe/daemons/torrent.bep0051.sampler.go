package daemons

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

const bep0051InfohashLen = 20

func NewSampler(q sqlx.Queryer, ttl time.Duration, cachepath string) *bep0051Sampler {
	return &bep0051Sampler{
		Q:    q,
		TTL:  ttl,
		Path: cachepath,
	}
}

type bep0051Sampler struct {
	Q    sqlx.Queryer
	TTL  time.Duration
	Path string
}

// snapshot returns the concatenated infohashes of public, seeding metadata.
// the result is cached on disk for TTL so repeated dht queries within the
// same window are served without hitting the database.
func (t bep0051Sampler) snapshot() (sample []byte, err error) {
	if info, serr := os.Stat(t.Path); serr == nil && time.Since(info.ModTime()) < t.TTL {
		if sample, err = os.ReadFile(t.Path); err == nil {
			return sample, nil
		}
	}

	var results []tracking.Metadata
	if err = sqlx.ScanInto(tracking.MetadataSamplePublic(context.Background(), t.Q), &results); err != nil {
		return nil, err
	}

	sample = make([]byte, 0, len(results)*bep0051InfohashLen)
	for _, m := range results {
		sample = append(sample, m.Infohash...)
	}

	if err = os.MkdirAll(filepath.Dir(t.Path), 0700); err != nil {
		return nil, err
	}

	if err = os.WriteFile(t.Path, sample, 0600); err != nil {
		return nil, err
	}

	return sample, nil
}

func (t bep0051Sampler) Snapshot(max int) (ttl uint, total uint, sample []byte) {
	ttl = uint(t.TTL / time.Second)

	full, err := t.snapshot()
	if err != nil {
		errorsx.Log(errorsx.Wrap(err, "unable to sample public metadata"))
		return ttl, 0, nil
	}

	total = uint(len(full) / bep0051InfohashLen)

	if limit := max * bep0051InfohashLen; max >= 0 && limit < len(full) {
		sample = full[:limit]
	} else {
		sample = full
	}

	return ttl, total, sample
}
