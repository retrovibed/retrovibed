package cmdmeta

type CloudBackup struct {
	Snapshot CloudBackupSnapshot `cmd:"" help:"upload an encrypted snapshot of the metadata database now"`
	Restore  CloudBackupRestore  `cmd:"" help:"restore the metadata database from the latest encrypted snapshot"`
}
