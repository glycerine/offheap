package offheap

import (
	"fmt"

	xxh64 "github.com/glycerine/xxhash-64"
)

// ByteKeyHashTable shows how to specialize HashTable to
// handle the []byte type as a key.
type ByteKeyHashTable HashTable

// xxHasher64 provides hashing for the byte-key (BK) interface to the hash table
//
// We use the 64-bit implimentation of XXHash for speed.
// see
//   https://github.com/OneOfOne/xxhash (github.com/glycerine/xxhash-64 version-locks)
//   http://fastcompression.blogspot.com/2014/07/xxhash-wider-64-bits.html
//
var xxHasher64 = xxh64.New64()

// NewByteKeyHashTable produces a new ByteKeyHashTable, one specialized for
// handling []byte as keys.
func NewByteKeyHashTable(initialSize uint64) *ByteKeyHashTable {
	return (*ByteKeyHashTable)(NewHashTable(initialSize))
}

// InsertBK is the insert function for []byte keys. By default only len(Key_t) bytes are used in the key.
func (t *ByteKeyHashTable) InsertBK(bytekey []byte, value interface{}) bool {
	xxHasher64.Reset()
	min := minimum(len(Key_t{}), len(bytekey))
	_, err := xxHasher64.Write(bytekey[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell, ok := (*HashTable)(t).Insert(hashkey)
	copy(cell.ByteKey[:], bytekey)
	cell.SetValue(value)
	return ok
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LookupBK is the lookup function for []byte keys. By default only len(Key_t) bytes are used in the key.
func (t *ByteKeyHashTable) LookupBK(bytekey []byte) (Val_t, bool) {
	xxHasher64.Reset()
	min := minimum(len(Key_t{}), len(bytekey))
	_, err := xxHasher64.Write(bytekey[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell := (*HashTable)(t).Lookup(hashkey)
	if cell == nil {
		return Val_t{}, false
	}
	return cell.Value, true
}

// DeleteBK removes an entry with a []byte key. By default only len(Key_t) bytes are used in the key.
func (t *ByteKeyHashTable) DeleteBK(bytekey []byte) bool {
	xxHasher64.Reset()
	min := minimum(len(Key_t{}), len(bytekey))
	_, err := xxHasher64.Write(bytekey[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell := (*HashTable)(t).Lookup(hashkey)
	if cell == nil {
		return false
	}

	(*HashTable)(t).DeleteCell(cell)
	return true
}

// StringHashTable shows how to specialize HashTable to
// handle strings as keys.
type StringHashTable HashTable

// NewStringHashTable produces a new StringHashTable, one specialized for
// handling keys of type string.
func NewStringHashTable(initialSize uint64) *StringHashTable {
	return (*StringHashTable)(NewHashTable(initialSize))
}

// InsertStringKey inserts a value with a key that is a string.
func (t *StringHashTable) InsertStringKey(strkey string, value interface{}) bool {
	xxHasher64.Reset()

	min := minimum(len(Key_t{}), len(strkey))
	var bytekey Key_t
	copy(bytekey[:], []byte(strkey))
	_, err := xxHasher64.Write(bytekey[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell, ok := (*HashTable)(t).Insert(hashkey)
	cell.ByteKey = bytekey
	cell.SetValue(value)
	//fmt.Printf("assigned value : '%v'  to key: '%v', with strkey: '%v'\n", value, hashkey, strkey)

	return ok
}

// LookupStringKey looks up a value based on a key that is a string.
func (t *StringHashTable) LookupStringKey(strkey string) (Val_t, bool) {
	xxHasher64.Reset()
	min := minimum(len(Key_t{}), len(strkey))
	var bytekey Key_t
	copy(bytekey[:], []byte(strkey))
	_, err := xxHasher64.Write(bytekey[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell := (*HashTable)(t).Lookup(hashkey)
	if cell == nil {
		return Val_t{}, false
	}
	return cell.Value, true
}

// DeleteStringKey deletes the cell (if there is one) that has
// been previously inserted with the given strkey string key.
func (t *StringHashTable) DeleteStringKey(strkey string) bool {
	xxHasher64.Reset()
	min := minimum(len(Key_t{}), len(strkey))
	var bytekey Key_t
	copy(bytekey[:], []byte(strkey))
	_, err := xxHasher64.Write(bytekey[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell := (*HashTable)(t).Lookup(hashkey)
	if cell == nil {
		return false
	}

	(*HashTable)(t).DeleteCell(cell)
	return true
}

// DumpStringKey provides a diagnostic printout of the
// contents of a hashtable that is using strings as keys.
func (t *StringHashTable) DumpStringKey() {

	fmt.Printf(" DumpStringKey(): (table ArraySize: %d\n", t.ArraySize)
	for it := NewIterator((*HashTable)(t)); it.Cur != nil; it.Next() {
		fmt.Printf("  '%v' -> %v\n", string(it.Cur.ByteKey[:]), it.Cur.Value)
	}

}
