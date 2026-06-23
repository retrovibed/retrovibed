package ddisc

// Snapshot is a point in time view of the discovery configuration and the
// locally assigned partition, intended for diagnostics reporting.
type Snapshot struct {
	Enabled        bool
	Ratio          uint32
	Partitions     uint32
	Workloads      uint32
	LocalPartition string
}
