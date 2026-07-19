package ddisc

import (
	"cmp"
	"iter"
	"math"
	"regexp"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// ErrNoCandidate indicates every Discovered row for a known-media-id was
// either absent or hard-rejected by the Policy.
var ErrNoCandidate = errorsx.String("no acceptable candidate")

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
	combined := uint16(bracket)*1000 + uint16(clampInt(customFormatScore, -500, 500)+500)

	// PolicyRank ranks toward zero (best), so invert combined - which
	// scores toward its max (best) - rather than subtracting from a
	// magic ceiling.
	d.PolicyRank = ^combined

	return nil
}

// Compare orders candidates by PolicyRank ascending, then Health descending,
// then Bytes descending: negative when a should be preferred over b.
func Compare(a, b Discovered) int {
	return cmp.Or(
		cmp.Compare(a.PolicyRank, b.PolicyRank),
		cmp.Compare(b.Health, a.Health),
		cmp.Compare(b.Bytes, a.Bytes),
	)
}

// Select reduces seq - a sequence of candidates already ranked by a Policy
// (e.g. by Discover, which ranks every candidate before yielding it) - to
// the single best non-rejected one. Returns ErrNoCandidate if seq yielded
// nothing or every candidate was rejected.
func Select(seq iter.Seq[Discovered]) (Discovered, error) {
	var (
		best  = Worst()
		found bool
	)

	for d := range seq {
		if d.PolicyRejection != "" {
			continue
		}

		if Compare(d, best) < 0 {
			best, found = d, true
		}
	}

	if !found {
		return best, ErrNoCandidate
	}

	return best, nil
}
