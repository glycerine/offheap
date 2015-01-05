package offheap

import (
	"github.com/glycerine/go-capnproto"
	kv "github.com/glycerine/go-offheap-hashtable/keyval"
)

// demonstrate use of custom key and value with
// the hashtable. The structs we use are in keyval/account.go.
// The AcctId struct acts as the key.
// The Account struct acts as the value.
//
// There are three functions to write: Insert(), Lookup(), and Delete().
//

func HashAcctId(acctid []byte) uint64 {
	xxHasher64.Reset()
	min := minimum(len(key_t{}), len(acctid))
	_, err := xxHasher64.Write(acctid[:min])
	if err != nil {
		panic(err)
	}
	hashkey := xxHasher64.Sum64()
	return hashkey
}

func (t *HashTable) InsertAcct(acctid string, value *kv.Account) bool {
	aid := []byte(acctid)
	hashkey := HashAcctId(aid)
	cell, ok := t.Insert(hashkey)
	copy(cell.ByteKey[:], aid)
	cell.Value = capn.Object(kv.AccountGoToCapn(&t.seg, value))
	return ok
}

func (t *HashTable) LookupAcct(acctid string) (kv.AccountCapn, bool) {
	aid := []byte(acctid)
	hashkey := HashAcctId(aid)
	cell := t.Lookup(hashkey)
	if cell == nil {
		return kv.AccountCapn{}, false
	}
	return kv.AccountCapn(cell.Value), true
}

func (t *HashTable) DeleteAcct(acctid string) bool {
	aid := []byte(acctid)
	hashkey := HashAcctId(aid)
	cell := t.Lookup(hashkey)
	if cell == nil {
		return false
	}

	t.DeleteCell(cell)
	return true
}
