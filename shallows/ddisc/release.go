package ddisc

import "regexp"

type ReleaseSource int

const (
	SourceUnknown ReleaseSource = iota
	SourceSDTV
	SourceHDTV
	SourceWEBRip
	SourceWEBDL
	SourceBluRay
	SourceRemux
)

type ReleaseResolution int

const (
	ResolutionUnknown ReleaseResolution = iota
	Resolution480p
	Resolution720p
	Resolution1080p
	Resolution2160p
)

type ReleaseInfo struct {
	Source     ReleaseSource
	Resolution ReleaseResolution
	HDR        bool // HDR10, HDR10+, DV/Dolby Vision
	Atmos      bool // Dolby Atmos / TrueHD Atmos
	Remux      bool
}

var (
	releaseRemux  = regexp.MustCompile(`(?i)\bremux\b`)
	releaseBluRay = regexp.MustCompile(`(?i)\b(blu-?ray|bdrip|brrip)\b`)
	releaseWEBDL  = regexp.MustCompile(`(?i)\bweb-?dl\b`)
	releaseWEBRip = regexp.MustCompile(`(?i)\bweb-?rip\b`)
	releaseHDTV   = regexp.MustCompile(`(?i)\bhdtv\b`)
	releaseSDTV   = regexp.MustCompile(`(?i)\b(sdtv|dvdrip)\b`)
	release2160p  = regexp.MustCompile(`(?i)\b(2160p|4k)\b`)
	release1080p  = regexp.MustCompile(`(?i)\b1080p\b`)
	release720p   = regexp.MustCompile(`(?i)\b720p\b`)
	release480p   = regexp.MustCompile(`(?i)\b480p\b`)
	releaseHDR    = regexp.MustCompile(`(?i)\b(hdr10\+?|dolby.?vision|dv)\b`)
	releaseAtmos  = regexp.MustCompile(`(?i)\batmos\b`)
)

// ExtractRelease parses a torrent/release title for the source, resolution,
// and format tags a ranking Policy uses to bracket quality. There is
// deliberately no attempt to be exhaustive - unmatched fields stay at their
// Unknown zero value.
func ExtractRelease(title string) (info ReleaseInfo) {
	switch {
	case releaseRemux.MatchString(title):
		info.Source = SourceRemux
		info.Remux = true
	case releaseBluRay.MatchString(title):
		info.Source = SourceBluRay
	case releaseWEBDL.MatchString(title):
		info.Source = SourceWEBDL
	case releaseWEBRip.MatchString(title):
		info.Source = SourceWEBRip
	case releaseHDTV.MatchString(title):
		info.Source = SourceHDTV
	case releaseSDTV.MatchString(title):
		info.Source = SourceSDTV
	default:
		info.Source = SourceUnknown
	}

	switch {
	case release2160p.MatchString(title):
		info.Resolution = Resolution2160p
	case release1080p.MatchString(title):
		info.Resolution = Resolution1080p
	case release720p.MatchString(title):
		info.Resolution = Resolution720p
	case release480p.MatchString(title):
		info.Resolution = Resolution480p
	default:
		info.Resolution = ResolutionUnknown
	}

	info.HDR = releaseHDR.MatchString(title)
	info.Atmos = releaseAtmos.MatchString(title)

	return info
}
