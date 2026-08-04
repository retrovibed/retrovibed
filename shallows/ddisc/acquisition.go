package ddisc

// AcquisitionState tracks the search/internal-database lifecycle of a
// Discovered candidate, independent of how the underlying content is
// treated (see ContentMime). Ordinals mirror ddisc.discovery.proto's
// AcquisitionState enum exactly - keep them in sync.
type AcquisitionState int32

const (
	AcquisitionStateUnknown     AcquisitionState = 0
	AcquisitionStateEphemeral   AcquisitionState = 1 // search result, not persisted.
	AcquisitionStateAvailable   AcquisitionState = 2 // persisted but download was not initiated; e.g. queued as a Recommendation.
	AcquisitionStateDownloading AcquisitionState = 3 // actively being downloaded.
	AcquisitionStateCompleted   AcquisitionState = 4 // download completed.
)
