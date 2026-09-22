package ddisc

import (
	"context"
	"log"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// DownloadDiscovered resolves d's torrent metadata via importer, persists d
// (correcting whatever placeholder infohash it carried until now), and
// either marks the resulting torrent for auto-download or records a
// Recommendation for later review. Full info-dict resolution (if not
// already available from d.URI) is left to the normal download/resume
// machinery (see daemons.ResumeDownloads) once initiated_at is set - same as
// every other ingestion path (RSS, Locate, etc.).
func DownloadDiscovered(ctx context.Context, db sqlx.Queryer, importer tracking.URIImport, d Discovered, acquisition AcquisitionState, options ...DiscoveredOption) (Discovered, tracking.Metadata, error) {
	lmd, err := importer.Import(
		ctx,
		d.URI,
		tracking.MetadataOptionKnownMediaID(d.KnownMediaID),
		tracking.MetadataOptionAutoDescription,
		tracking.MetadataOptionAutoHidden,
	)
	if err != nil {
		return d, lmd, errorsx.Wrapf(err, "unable to import uri for download %s", d.ID)
	}

	// import resolved the real infohash (parsed from the magnet, or hashed
	// from the actually-fetched .torrent bytes) - persist d now, correcting
	// whatever placeholder infohash it carried until this point. It also
	// clears d's default-private flag (see NewDiscoveredFromImport) to
	// whatever the fetched torrent's info dict actually says - the only
	// point BEP 27 privacy becomes known.
	d = langx.Clone(
		d,
		DiscoveredOptionAcquisitionState(acquisition),
		DiscoveredOptionInfoHash(lmd.Infohash),
		DiscoveredOptionPrivate(lmd.Private),
		langx.Compose(options...),
	)
	if err = DiscoveredInsertWithDefaults(ctx, db, d).Scan(&d); err != nil {
		return d, lmd, errorsx.Wrapf(err, "unable to persist discovered candidate %s", d.ID)
	}

	switch acquisition {
	case AcquisitionStateDownloading:
		if err = tracking.MetadataAutoDownloadByID(ctx, db, lmd.ID).Scan(&lmd); err != nil {
			return d, lmd, errorsx.Wrapf(err, "unable to mark torrent for download from uri %s", d.ID)
		}
		log.Println("marked for download", lmd.ID, lmd.Description)
	default:
		var rec library.Recommendation
		if err = library.RecommendationInsertWithDefaults(ctx, db, RecommendationFromDiscovered(d)).Scan(&rec); err != nil {
			return d, lmd, errorsx.Wrap(err, "unable to record recommendation for discovered media")
		}
		log.Println("recommendation created", rec.ID, rec.Title)
	}

	return d, lmd, nil
}
