package blockcache

// FileIndex is a test-only function to extract the current atomic index from a File.
// It is exposed only when the "test" build tag is active.
func FileIndex(f *File) uint64 {
	return f.index.Load()
}
