package offheap

// AUTO GENERATED - DO NOT EDIT

import (
	"unsafe"

	C "github.com/glycerine/go-capnproto"
)

type CellCapn C.Struct

func NewCellCapn(s *C.Segment) CellCapn      { return CellCapn(s.NewStruct(8, 2)) }
func NewRootCellCapn(s *C.Segment) CellCapn  { return CellCapn(s.NewRootStruct(8, 2)) }
func AutoNewCellCapn(s *C.Segment) CellCapn  { return CellCapn(s.NewStructAR(8, 2)) }
func ReadRootCellCapn(s *C.Segment) CellCapn { return CellCapn(s.Root(0).ToStruct()) }
func (s CellCapn) UnHashedKey() uint64       { return C.Struct(s).Get64(0) }
func (s CellCapn) SetUnHashedKey(v uint64)   { C.Struct(s).Set64(0, v) }
func (s CellCapn) ByteKey() []byte           { return C.Struct(s).GetObject(0).ToData() }
func (s CellCapn) SetByteKey(v []byte)       { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }
func (s CellCapn) Value() []byte             { return C.Struct(s).GetObject(1).ToData() }
func (s CellCapn) SetValue(v []byte)         { C.Struct(s).SetObject(1, s.Segment.NewData(v)) }

func (s CellCapn) MarshalJSON() (bs []byte, err error) { return }

type CellCapn_List C.PointerList

func NewCellCapnList(s *C.Segment, sz int) CellCapn_List {
	return CellCapn_List(s.NewCompositeList(8, 2, sz))
}
func (s CellCapn_List) Len() int          { return C.PointerList(s).Len() }
func (s CellCapn_List) At(i int) CellCapn { return CellCapn(C.PointerList(s).At(i).ToStruct()) }
func (s CellCapn_List) ToArray() []CellCapn {
	return *(*[]CellCapn)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
func (s CellCapn_List) Set(i int, item CellCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type HashTableCapn C.Struct

func NewHashTableCapn(s *C.Segment) HashTableCapn      { return HashTableCapn(s.NewStruct(40, 3)) }
func NewRootHashTableCapn(s *C.Segment) HashTableCapn  { return HashTableCapn(s.NewRootStruct(40, 3)) }
func AutoNewHashTableCapn(s *C.Segment) HashTableCapn  { return HashTableCapn(s.NewStructAR(40, 3)) }
func ReadRootHashTableCapn(s *C.Segment) HashTableCapn { return HashTableCapn(s.Root(0).ToStruct()) }
func (s HashTableCapn) Cells() uint64                  { return C.Struct(s).Get64(0) }
func (s HashTableCapn) SetCells(v uint64)              { C.Struct(s).Set64(0, v) }
func (s HashTableCapn) CellSz() uint64                 { return C.Struct(s).Get64(8) }
func (s HashTableCapn) SetCellSz(v uint64)             { C.Struct(s).Set64(8, v) }
func (s HashTableCapn) ArraySize() uint64              { return C.Struct(s).Get64(16) }
func (s HashTableCapn) SetArraySize(v uint64)          { C.Struct(s).Set64(16, v) }
func (s HashTableCapn) Population() uint64             { return C.Struct(s).Get64(24) }
func (s HashTableCapn) SetPopulation(v uint64)         { C.Struct(s).Set64(24, v) }
func (s HashTableCapn) ZeroUsed() bool                 { return C.Struct(s).Get1(256) }
func (s HashTableCapn) SetZeroUsed(v bool)             { C.Struct(s).Set1(256, v) }
func (s HashTableCapn) ZeroCell() CellCapn             { return CellCapn(C.Struct(s).GetObject(0).ToStruct()) }
func (s HashTableCapn) SetZeroCell(v CellCapn)         { C.Struct(s).SetObject(0, C.Object(v)) }
func (s HashTableCapn) Offheap() []byte                { return C.Struct(s).GetObject(1).ToData() }
func (s HashTableCapn) SetOffheap(v []byte)            { C.Struct(s).SetObject(1, s.Segment.NewData(v)) }
func (s HashTableCapn) Mmm() MmapMallocCapn {
	return MmapMallocCapn(C.Struct(s).GetObject(2).ToStruct())
}
func (s HashTableCapn) SetMmm(v MmapMallocCapn) { C.Struct(s).SetObject(2, C.Object(v)) }

func (s HashTableCapn) MarshalJSON() (bs []byte, err error) { return }

type HashTableCapn_List C.PointerList

func NewHashTableCapnList(s *C.Segment, sz int) HashTableCapn_List {
	return HashTableCapn_List(s.NewCompositeList(40, 3, sz))
}
func (s HashTableCapn_List) Len() int { return C.PointerList(s).Len() }
func (s HashTableCapn_List) At(i int) HashTableCapn {
	return HashTableCapn(C.PointerList(s).At(i).ToStruct())
}
func (s HashTableCapn_List) ToArray() []HashTableCapn {
	return *(*[]HashTableCapn)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
func (s HashTableCapn_List) Set(i int, item HashTableCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type IteratorCapn C.Struct

func NewIteratorCapn(s *C.Segment) IteratorCapn      { return IteratorCapn(s.NewStruct(8, 2)) }
func NewRootIteratorCapn(s *C.Segment) IteratorCapn  { return IteratorCapn(s.NewRootStruct(8, 2)) }
func AutoNewIteratorCapn(s *C.Segment) IteratorCapn  { return IteratorCapn(s.NewStructAR(8, 2)) }
func ReadRootIteratorCapn(s *C.Segment) IteratorCapn { return IteratorCapn(s.Root(0).ToStruct()) }
func (s IteratorCapn) Tab() HashTableCapn            { return HashTableCapn(C.Struct(s).GetObject(0).ToStruct()) }
func (s IteratorCapn) SetTab(v HashTableCapn)        { C.Struct(s).SetObject(0, C.Object(v)) }
func (s IteratorCapn) Pos() int64                    { return int64(C.Struct(s).Get64(0)) }
func (s IteratorCapn) SetPos(v int64)                { C.Struct(s).Set64(0, uint64(v)) }
func (s IteratorCapn) Cur() CellCapn                 { return CellCapn(C.Struct(s).GetObject(1).ToStruct()) }
func (s IteratorCapn) SetCur(v CellCapn)             { C.Struct(s).SetObject(1, C.Object(v)) }

func (s IteratorCapn) MarshalJSON() (bs []byte, err error) { return }

type IteratorCapn_List C.PointerList

func NewIteratorCapnList(s *C.Segment, sz int) IteratorCapn_List {
	return IteratorCapn_List(s.NewCompositeList(8, 2, sz))
}
func (s IteratorCapn_List) Len() int { return C.PointerList(s).Len() }
func (s IteratorCapn_List) At(i int) IteratorCapn {
	return IteratorCapn(C.PointerList(s).At(i).ToStruct())
}
func (s IteratorCapn_List) ToArray() []IteratorCapn {
	return *(*[]IteratorCapn)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
func (s IteratorCapn_List) Set(i int, item IteratorCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type MmapMallocCapn C.Struct

func NewMmapMallocCapn(s *C.Segment) MmapMallocCapn      { return MmapMallocCapn(s.NewStruct(24, 2)) }
func NewRootMmapMallocCapn(s *C.Segment) MmapMallocCapn  { return MmapMallocCapn(s.NewRootStruct(24, 2)) }
func AutoNewMmapMallocCapn(s *C.Segment) MmapMallocCapn  { return MmapMallocCapn(s.NewStructAR(24, 2)) }
func ReadRootMmapMallocCapn(s *C.Segment) MmapMallocCapn { return MmapMallocCapn(s.Root(0).ToStruct()) }
func (s MmapMallocCapn) Path() string                    { return C.Struct(s).GetObject(0).ToText() }
func (s MmapMallocCapn) SetPath(v string)                { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s MmapMallocCapn) Fd() int64                       { return int64(C.Struct(s).Get64(0)) }
func (s MmapMallocCapn) SetFd(v int64)                   { C.Struct(s).Set64(0, uint64(v)) }
func (s MmapMallocCapn) FileBytesLen() int64             { return int64(C.Struct(s).Get64(8)) }
func (s MmapMallocCapn) SetFileBytesLen(v int64)         { C.Struct(s).Set64(8, uint64(v)) }
func (s MmapMallocCapn) BytesAlloc() int64               { return int64(C.Struct(s).Get64(16)) }
func (s MmapMallocCapn) SetBytesAlloc(v int64)           { C.Struct(s).Set64(16, uint64(v)) }
func (s MmapMallocCapn) Mem() []byte                     { return C.Struct(s).GetObject(1).ToData() }
func (s MmapMallocCapn) SetMem(v []byte)                 { C.Struct(s).SetObject(1, s.Segment.NewData(v)) }

func (s MmapMallocCapn) MarshalJSON() (bs []byte, err error) { return }

type MmapMallocCapn_List C.PointerList

func NewMmapMallocCapnList(s *C.Segment, sz int) MmapMallocCapn_List {
	return MmapMallocCapn_List(s.NewCompositeList(24, 2, sz))
}
func (s MmapMallocCapn_List) Len() int { return C.PointerList(s).Len() }
func (s MmapMallocCapn_List) At(i int) MmapMallocCapn {
	return MmapMallocCapn(C.PointerList(s).At(i).ToStruct())
}
func (s MmapMallocCapn_List) ToArray() []MmapMallocCapn {
	return *(*[]MmapMallocCapn)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
func (s MmapMallocCapn_List) Set(i int, item MmapMallocCapn) { C.PointerList(s).Set(i, C.Object(item)) }
