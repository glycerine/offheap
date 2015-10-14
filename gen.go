//go:generate genny -pkg=offheap -in=template/offheap.go -out=somestruct_offheap_gen.go gen "_T_=SomeStruct _LT_=SomeStruct _C_=EmptyStructSomeStruct"
//go:generate genny -pkg=offheap -in=template/offheap.go -out=int_offheap_gen.go gen "_T_=int _LT_=Int _C_=CustomMetadataIntTest"
package offheap

type CustomMetadataIntTest struct {
	A int
	B float64
	C [10]byte
}

func (h *HashTableInt) InsertInt(k uint64, v int) {
	c, _ := h.Insert(k)
	c.Value = v
}

func (h *HashTableInt) LookupInt(key uint64) (int, bool) {
	c := h.Lookup(key)
	if c != nil {
		return c.Value, true
	}
	return 0, false
}
