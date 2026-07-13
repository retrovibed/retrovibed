package ddisc

import (
	"context"
	"math"
	"regexp"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// ErrNoCandidate indicates every Discovered row for a known-media-id was
// either absent or hard-rejected by the Policy.
var ErrNoCandidate = errorsx.String("no acceptable candidate")

// healthMidTierDefault is used when a candidate's Health is unknown (its
// zero value) rather than actually reported as zero seeds - e.g. anything
// not sourced from a search plugin, since nothing else currently reports a
// real health number. A placeholder pending real DHT BEP33 scrape support.
const healthMidTierDefault uint32 = 25

// Policy ranks a discovered candidate for download selection.
type Policy interface {
	// Rank sets d.Health, d.PolicyRank, and d.PolicyRejection in place.
	// A non-nil error is a genuine failure (e.g. a bug in the ranking
	// computation), not a rejection - rejections are recorded as data on
	// d (PolicyRejection set to the matched dealbreaker, PolicyRank set
	// to the worst-possible sentinel), not surfaced as a Go error.
	Rank(d *Discovered) error
}

// rejectTerms are release-title dealbreakers modeled on Radarr/Sonarr
// defaults: cam/telesync-class releases are rejected outright rather than
// merely down-ranked.
var rejectTerms = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bhdcam\b`),
	regexp.MustCompile(`(?i)\bcam\b`),
	regexp.MustCompile(`(?i)\btelesync\b`),
	regexp.MustCompile(`(?i)\bts\b`),
	regexp.MustCompile(`(?i)\btelecine\b`),
	regexp.MustCompile(`(?i)\bworkprint\b`),
}

var screenerTerm = regexp.MustCompile(`(?i)\bscreener\b`)

func sourceRank(s ReleaseSource) int {
	switch s {
	case SourceRemux:
		return 6
	case SourceBluRay:
		return 5
	case SourceWEBDL:
		return 4
	case SourceWEBRip:
		return 3
	case SourceHDTV:
		return 2
	case SourceSDTV:
		return 1
	default:
		return 0
	}
}

func resolutionRank(r ReleaseResolution) int {
	switch r {
	case Resolution2160p:
		return 4
	case Resolution1080p:
		return 3
	case Resolution720p:
		return 2
	case Resolution480p:
		return 1
	default:
		return 0
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

type defaultPolicy struct{}

// DefaultPolicy is a Radarr/Sonarr-conventions ranking Policy: hard-reject
// CAM/TS-class releases, prefer higher quality brackets (resolution first,
// then source), break ties with custom-format scoring (HDR/Atmos/Remux
// bonuses, Screener penalty), then swarm health, then file size.
func DefaultPolicy() Policy {
	return defaultPolicy{}
}

func (defaultPolicy) Rank(d *Discovered) error {
	for _, term := range rejectTerms {
		if term.MatchString(d.Title) {
			d.PolicyRejection = term.String()
			d.PolicyRank = math.MaxUint16
			return nil
		}
	}

	release := ExtractRelease(d.Title)

	bracket := resolutionRank(release.Resolution)*10 + sourceRank(release.Source)

	customFormatScore := 0
	if release.HDR {
		customFormatScore += 100
	}
	if release.Atmos {
		customFormatScore += 30
	}
	if release.Remux {
		customFormatScore += 50
	}
	if screenerTerm.MatchString(d.Title) {
		customFormatScore -= 50
	}

	// bracket must strictly dominate customFormatScore: even the best
	// possible custom-format score can't outweigh a full bracket step.
	combined := uint32(bracket)*1000 + uint32(clampInt(customFormatScore, -500, 500)+500)
	d.PolicyRank = uint16(65534 - combined)

	if d.Health == 0 {
		d.Health = healthMidTierDefault
	}

	return nil
}

// RankAndSelect ranks every Discovered row for knownMediaID with policy,
// persists the result of each (health/policy_rank/policy_rejection), and
// returns the lowest-PolicyRank non-rejected candidate. Health and Bytes
// break ties between equally-ranked candidates.
func RankAndSelect(ctx context.Context, q sqlx.Queryer, policy Policy, knownMediaID string) (Discovered, error) {
	s := sqlx.Scan(DiscoveredByKnownID(ctx, q, knownMediaID))

	var (
		best  Discovered
		found bool
	)

	for d := range s.Iter() {
		if err := policy.Rank(&d); err != nil {
			return best, err
		}

		if err := DiscoveredRank(ctx, q, d.ID, d.Health, d.PolicyRank, d.PolicyRejection).Scan(&d); err != nil {
			return best, err
		}

		if d.PolicyRejection != "" {
			continue
		}

		if !found ||
			d.PolicyRank < best.PolicyRank ||
			(d.PolicyRank == best.PolicyRank && (d.Health > best.Health ||
				(d.Health == best.Health && d.Bytes > best.Bytes))) {
			best, found = d, true
		}
	}

	if err := s.Err(); err != nil {
		return best, err
	}

	if !found {
		return best, ErrNoCandidate
	}

	return best, nil
}
