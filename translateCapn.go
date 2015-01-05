package offheap

import (
	"fmt"
	"io"

	capn "github.com/glycerine/go-capnproto"
)

func (s *Cell) Save(w io.Writer) {
	seg := capn.NewBuffer(nil)
	CellGoToCapn(seg, s)
	seg.WriteTo(w)
}

func (s *Cell) Load(r io.Reader) {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
	}
	z := ReadRootCellCapn(capMsg)
	CellCapnToGo(z, s)
}

func CellCapnToGo(src CellCapn, dest *Cell) *Cell {
	if dest == nil {
		dest = &Cell{}
	}
	dest.UnHashedKey = uint64(src.UnHashedKey())

	return dest
}

func CellGoToCapn(seg *capn.Segment, src *Cell) CellCapn {
	dest := AutoNewCellCapn(seg)
	dest.SetUnHashedKey(uint64(src.UnHashedKey))

	return dest
}

func (s *HashTable) Save(w io.Writer) {
	seg := capn.NewBuffer(nil)
	HashTableGoToCapn(seg, s)
	seg.WriteTo(w)
}

func (s *HashTable) Load(r io.Reader) {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
	}
	z := ReadRootHashTableCapn(capMsg)
	HashTableCapnToGo(z, s)
}

func HashTableCapnToGo(src HashTableCapn, dest *HashTable) *HashTable {
	if dest == nil {
		dest = &HashTable{}
	}
	dest.CellSz = uint64(src.CellSz())
	dest.ArraySize = uint64(src.ArraySize())
	dest.Population = uint64(src.Population())

	var n int

	// Offheap
	n = len(src.Offheap())
	dest.Offheap = make([]byte, n)
	for i := 0; i < n; i++ {
		dest.Offheap[i] = byte(src.Offheap()[i])
	}

	return dest
}

func HashTableGoToCapn(seg *capn.Segment, src *HashTable) HashTableCapn {
	dest := AutoNewHashTableCapn(seg)
	dest.SetCellSz(uint64(src.CellSz))
	dest.SetArraySize(uint64(src.ArraySize))
	dest.SetPopulation(uint64(src.Population))

	mylist1 := seg.NewUInt8List(len(src.Offheap))
	for i := range src.Offheap {
		mylist1.Set(i, uint8(src.Offheap[i]))
	}
	dest.SetOffheap(mylist1.ToArray())

	return dest
}

func (s *Iterator) Save(w io.Writer) {
	seg := capn.NewBuffer(nil)
	IteratorGoToCapn(seg, s)
	seg.WriteTo(w)
}

func (s *Iterator) Load(r io.Reader) {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
	}
	z := ReadRootIteratorCapn(capMsg)
	IteratorCapnToGo(z, s)
}

func IteratorCapnToGo(src IteratorCapn, dest *Iterator) *Iterator {
	if dest == nil {
		dest = &Iterator{}
	}
	dest.Pos = int64(src.Pos())

	return dest
}

func IteratorGoToCapn(seg *capn.Segment, src *Iterator) IteratorCapn {
	dest := AutoNewIteratorCapn(seg)
	dest.SetPos(src.Pos)

	return dest
}

func (s *MmapMalloc) Save(w io.Writer) {
	seg := capn.NewBuffer(nil)
	MmapMallocGoToCapn(seg, s)
	seg.WriteTo(w)
}

func (s *MmapMalloc) Load(r io.Reader) {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
	}
	z := ReadRootMmapMallocCapn(capMsg)
	MmapMallocCapnToGo(z, s)
}

func MmapMallocCapnToGo(src MmapMallocCapn, dest *MmapMalloc) *MmapMalloc {
	if dest == nil {
		dest = &MmapMalloc{}
	}
	dest.Path = src.Path()
	dest.Fd = int(src.Fd())
	dest.FileBytesLen = int64(src.FileBytesLen())
	dest.BytesAlloc = int64(src.BytesAlloc())

	var n int

	// Mem
	n = len(src.Mem())
	dest.Mem = make([]byte, n)
	for i := 0; i < n; i++ {
		dest.Mem[i] = byte(src.Mem()[i])
	}

	return dest
}

func MmapMallocGoToCapn(seg *capn.Segment, src *MmapMalloc) MmapMallocCapn {
	dest := AutoNewMmapMallocCapn(seg)
	dest.SetPath(src.Path)
	dest.SetFd(int64(src.Fd))
	dest.SetFileBytesLen(src.FileBytesLen)
	dest.SetBytesAlloc(src.BytesAlloc)

	mylist1 := seg.NewUInt8List(len(src.Mem))
	for i := range src.Mem {
		mylist1.Set(i, uint8(src.Mem[i]))
	}
	dest.SetMem(mylist1.ToArray())

	return dest
}

func SliceByteToUInt8List(seg *capn.Segment, m []byte) capn.UInt8List {
	lst := seg.NewUInt8List(len(m))
	for i := range m {
		lst.Set(i, uint8(m[i]))
	}
	return lst
}

func UInt8ListToSliceByte(p capn.UInt8List) []byte {
	v := make([]byte, p.Len())
	for i := range v {
		v[i] = byte(p.At(i))
	}
	return v
}
