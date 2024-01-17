package sst

type SstEntry struct {
	Key     string
	Value   []byte
	Deleted bool
}
