package lsm

import (
	"github.com/google/btree"
)

func New(path string, bufSize int) *LsmTree {
	return &LsmTree{
		memtbl: btree.New(2),
	}
}

// ResetDB removes all data from disk
func (tree *LsmTree) ResetDB() {
}

// Set will add (or update) an entry in the tree with the corresponding key/value.
func (tree *LsmTree) Set(k string, value []byte) {
	tree.set(k, value, false)
}

// Delete will remove the corresponding key from the tree.
// Note the actual key/value may not be removed from memory or disk immediately.
// One or more merge/compact must run before data is removed from disk.
func (tree *LsmTree) Delete(k string) {
	var val []byte
	tree.set(k, val, true)
}

// Increment will add one to the integer counter specified by the given key,
// and the most recent value will be returned.
// New counters return a value of 0.
func (tree *LsmTree) Increment(k string) uint32 {
	var result uint32

	return result
}

// Get looks up the given key and returns a boolean indicating if a value was
// found and (if found) the corresponding value as a byte array.
func (tree *LsmTree) Get(k string) ([]byte, bool) {
	val, ok := tree.get(k)
	return val, ok
}

type Item struct {
	Key     string
	Value   []byte
	Deleted bool
}

func (i *Item) Less(than btree.Item) bool {
	return i.Key < than.(*Item).Key
}

func (tree *LsmTree) setInMemtbl(k string, value []byte, deleted bool) {
	tree.lock.Lock()
	defer tree.lock.Unlock()
	tree.memtbl.ReplaceOrInsert(&Item{Key: k, Value: value, Deleted: deleted})
}
func (tree *LsmTree) set(k string, value []byte, deleted bool) {
	tree.setInMemtbl(k, value, deleted)
}

func (tree *LsmTree) get(k string) ([]byte, bool) {
	item := tree.memtbl.Get(&Item{Key: k})
	if item != nil {
		i := item.(*Item)
		if i.Deleted {
			return nil, false
		}
		return i.Value, true
	}
	return nil, false
}
