package offheap

import xxh64 "github.com/glycerine/xxhash-64"

// the byte-key (BK) interface to the hash table

// use the 64-bit implimentation of XXHash for speed.
// see
//   https://github.com/OneOfOne/xxhash (github.com/glycerine/xxhash-64 version-locks)
//   http://fastcompression.blogspot.com/2014/07/xxhash-wider-64-bits.html
//
var xxHasher64 = xxh64.New64()

func (t *HashTable) InsertBK(bytekey []byte, value interface{}) bool {
	xxHasher64.Reset()
	_, err := xxHasher64.Write(bytekey)
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell, ok := t.Insert(hashkey)
	cell.ByteKey = bytekey
	cell.Value = value
	return ok
}

func (t *HashTable) LookupBK(bytekey []byte) (interface{}, bool) {
	xxHasher64.Reset()
	_, err := xxHasher64.Write(bytekey)
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell := t.Lookup(hashkey)
	if cell == nil {
		return nil, false
	}
	return cell.Value, true
}

func (t *HashTable) InsertStringKey(strkey string, value interface{}) bool {
	xxHasher64.Reset()
	bytekey := []byte(strkey)
	_, err := xxHasher64.Write(bytekey)
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell, ok := t.Insert(hashkey)
	cell.ByteKey = bytekey
	cell.Value = value
	return ok
}

func (t *HashTable) LookupStringKey(strkey string) (interface{}, bool) {
	xxHasher64.Reset()
	bytekey := []byte(strkey)
	_, err := xxHasher64.Write(bytekey)
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	cell := t.Lookup(hashkey)
	if cell == nil {
		return nil, false
	}
	return cell.Value, true
}
