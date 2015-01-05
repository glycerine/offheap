package offheap

import (
	"fmt"
	"unsafe"

	capn "github.com/glycerine/go-capnproto"
)

// Copyright (C) 2015 by Jason E. Aten, Ph.D.
//

type key_t [64]byte

type Val_t [56]byte

type Cell struct {
	UnHashedKey uint64
	ByteKey     key_t
	Value       capn.Object
}

/*
func (cell *Cell) SetValue(v interface{}) {
	switch a := v.(type) {
	case string:
		cell.SetString(a)
	case int:
		cell.SetInt(a)
	default:
		panic("unsupported type")
	}
}
*/

func (cell *Cell) ZeroValue() {
	cell.Value = capn.Object{}
}

/*
func (cell *Cell) SetString(s string) {
	copy(cell.Value[:], []byte(s))
}
func (cell *Cell) GetString() string {
	return string([]byte(cell.Value[:]))
}

func (cell *Cell) SetInt(n int) {
	binary.LittleEndian.PutUint64(cell.Value[:8], uint64(n))
}
func (cell *Cell) GetInt() int {
	return int(binary.LittleEndian.Uint64(cell.Value[:8]))
}

func (v *Val_t) SetInt(n int) {
	binary.LittleEndian.PutUint64((*v)[:8], uint64(n))
}
func (v *Val_t) GetInt() int {
	return int(binary.LittleEndian.Uint64((*v)[:8]))
}

func (v *Val_t) SetString(s string) {
	copy((*v)[:], []byte(s))
}
func (v *Val_t) GetString() string {
	return string([]byte((*v)[:]))
}
*/

// The purpose here is to have hash table that can work away
// from Go's Garbage Collector, to avoid long GC pause times.
//
// We accomplish this by writing our own Malloc() and Free() implementation
// (see malloc.go) which requests memory directly from the OS.
// The keys, values, and entire hash table is kept on off-heap
// storage. This storage can also optionally be backed by memory mapped file
// for speedy persistence and fast startup times.
//
// Initial HashTable implementation inspired by the public domain C++ code of
//    https://github.com/preshing/CompareIntegerMaps
// See also
//    http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/
// for performance studies of the C++ code.
//
//
//----------------------------------------------
//  HashTable
//
//  Maps pointer-sized integers to Cell structures, which in turn hold Val_t
//   as well as key_t structures.
//
//  Uses open addressing with linear probing. This makes it very cache
//   friendly and thus very fast.
//
//  In the t.Cells array, UnHashedKey = 0 is reserved to indicate an unused
//   cell.
//
//  Actual value for key 0 (if any) is stored in t.ZeroCell.
//  The hash table automatically doubles in size when it becomes 75% full.
//  The hash table never shrinks in size, even after Clear(), unless you
//  explicitly call Compact().
//----------------------------------------------
//
type HashTable struct {
	Cells      uintptr      `capid:"0"`
	CellSz     uint64       `capid:"1"`
	ArraySize  uint64       `capid:"2"`
	Population uint64       `capid:"3"`
	ZeroUsed   bool         `capid:"4"`
	ZeroCell   Cell         `capid:"5"`
	Offheap    []byte       `capid:"6"`
	Mmm        MmapMalloc   `capid:"7"`
	seg        capn.Segment `capid:"skip"`
}

func NewHashTable(initialSize uint64) *HashTable {

	t := HashTable{
		CellSz: uint64(unsafe.Sizeof(Cell{})),
	}

	// off-heap and off-gc version
	t.ArraySize = initialSize
	t.Mmm = *Malloc(int64(t.ArraySize*t.CellSz), "")
	t.Offheap = t.Mmm.Mem
	t.Cells = (uintptr)(unsafe.Pointer(&t.Offheap[0]))

	t.seg = *capn.NewBuffer(t.Offheap)

	// off-gc but still on-heap version
	//	t.ArraySize = initialSize
	//	t.Offheap = make([]byte, t.ArraySize*t.CellSz)
	//	t.Cells = (uintptr)(unsafe.Pointer(&t.Offheap[0]))

	// on-heap version:
	//Cells:     make([]Cell, initialSize),

	return &t
}

// t.CellAt(pos); replaces t.Cells[pos]
func (t *HashTable) CellAt(pos uint64) *Cell {

	// off heap version
	return (*Cell)(unsafe.Pointer(uintptr(t.Cells) + uintptr(pos*t.CellSz)))

	// on heap version, back when t.Cells was []Cell
	//return &(t.Cells[pos])
}

func (t *HashTable) DestroyHashTable() {
	t.Mmm.Free()
}

// Basic operations

// return nil if not found
func (t *HashTable) Lookup(key uint64) *Cell {

	var cell *Cell

	if key == 0 {
		if t.ZeroUsed {
			return &t.ZeroCell
		}
		return nil

	} else {

		h := integerHash(uint64(key)) % t.ArraySize

		for {
			cell = t.CellAt(h)
			if cell.UnHashedKey == key {
				return cell
			}
			if cell.UnHashedKey == 0 {
				return nil
			}
			h++
			if h == t.ArraySize {
				h = 0
			}
		}
	}
}

// 2nd return value is false if already existed (and thus took no action)
func (t *HashTable) Insert(key uint64) (*Cell, bool) {

	VPrintf("\n ---- Insert(%v) called with t = \n", key)
	VDump(t)

	defer func() {
		VPrintf("\n ---- Insert(%v) done, with t = \n", key)
		VDump(t)
	}()

	var cell *Cell

	if key != 0 {

		for {
			h := integerHash(uint64(key)) % t.ArraySize

			for {
				cell = t.CellAt(h)
				if cell.UnHashedKey == key {
					// already exists
					return cell, false
				}
				if cell.UnHashedKey == 0 {
					if (t.Population+1)*4 >= t.ArraySize*3 {
						VPrintf("detected (t.Population+1)*4 >= t.ArraySize*3, i.e. %v >= %v, calling Repop with double the size\n", (t.Population+1)*4, t.ArraySize*3)
						t.Repopulate(t.ArraySize * 2)
						// resized, so start all over
						break
					}
					t.Population++
					cell.UnHashedKey = key
					return cell, true
				}

				h++
				if h == t.ArraySize {
					h = 0
				}

			}
		}
	} else {

		wasNew := false
		if !t.ZeroUsed {
			wasNew = true
			t.ZeroUsed = true
			t.Population++
			if t.Population*4 >= t.ArraySize*3 {

				t.Repopulate(t.ArraySize * 2)
			}
		}
		return &t.ZeroCell, wasNew
	}

}

func (t *HashTable) DeleteCell(cell *Cell) {

	if cell == &t.ZeroCell {
		// Delete zero cell
		if !t.ZeroUsed {
			panic("deleting zero element when not used")
		}
		t.ZeroUsed = false
		cell.ZeroValue()
		t.Population--
		return

	} else {

		pos := uint64((uintptr(unsafe.Pointer(cell)) - uintptr(unsafe.Pointer(t.Cells))) / uintptr(unsafe.Sizeof(Cell{})))

		// Delete from regular Cells
		if pos < 0 || pos >= t.ArraySize {
			panic(fmt.Sprintf("cell out of bounds: pos %v was < 0 or >= t.ArraySize == %v", pos, t.ArraySize))
		}

		if t.CellAt(pos).UnHashedKey == 0 {
			panic("zero UnHashedKey in non-zero Cell!")
		}

		// Remove this cell by shuffling neighboring Cells so there are no gaps in anyone's probe chain
		nei := pos + 1
		if nei >= t.ArraySize {
			nei = 0
		}
		var neighbor *Cell
		var circular_offset_ideal_pos int64
		var circular_offset_ideal_nei int64
		var cellPos *Cell

		for {
			neighbor = t.CellAt(nei)

			if neighbor.UnHashedKey == 0 {
				// There's nobody to swap with. Go ahead and clear this cell, then return
				cellPos = t.CellAt(pos)
				cellPos.UnHashedKey = 0
				cellPos.ZeroValue()
				t.Population--
				return
			}

			ideal := integerHash(neighbor.UnHashedKey) % t.ArraySize

			if pos >= ideal {
				circular_offset_ideal_pos = int64(pos) - int64(ideal)
			} else {
				// pos < ideal, so pos - ideal is negative, wrap-around has happened.
				circular_offset_ideal_pos = int64(t.ArraySize) - int64(ideal) + int64(pos)
			}

			if nei >= ideal {
				circular_offset_ideal_nei = int64(nei) - int64(ideal)
			} else {
				// nei < ideal, so nei - ideal is negative, wrap-around has happened.
				circular_offset_ideal_nei = int64(t.ArraySize) - int64(ideal) + int64(nei)
			}

			if circular_offset_ideal_pos < circular_offset_ideal_nei {
				// Swap with neighbor, then make neighbor the new cell to remove.
				*t.CellAt(pos) = *neighbor
				pos = nei
			}

			nei++
			if nei >= t.ArraySize {
				nei = 0
			}
		}
	}

}

func (t *HashTable) Clear() {
	// (Does not resize the array)
	// Clear regular Cells

	for i := range t.Offheap {
		t.Offheap[i] = 0
	}
	t.Population = 0

	// Clear zero cell
	t.ZeroUsed = false
	t.ZeroCell.ZeroValue()
}

func (t *HashTable) Compact() {
	t.Repopulate(Upper_power_of_two((t.Population*4 + 3) / 3))
}

func (t *HashTable) DeleteKey(key uint64) {
	value := t.Lookup(key)
	if value != nil {
		t.DeleteCell(value)
	}
}

func (t *HashTable) Repopulate(desiredSize uint64) {

	VPrintf("\n ---- Repopulate called with t = \n")
	VDump(t)

	if desiredSize&(desiredSize-1) != 0 {
		panic("desired size must be a power of 2")
	}
	if t.Population*4 > desiredSize*3 {
		panic("must have t.Population * 4  <= desiredSize * 3")
	}

	// Allocate new table
	s := NewHashTable(desiredSize)

	s.ZeroUsed = t.ZeroUsed
	if t.ZeroUsed {
		s.ZeroCell = t.ZeroCell
		s.Population++
	}

	// Iterate through old table t, copy into new table s.
	var c *Cell

	for i := uint64(0); i < t.ArraySize; i++ {
		c = t.CellAt(i)
		VPrintf("\n in oldCell copy loop, at i = %v, and c = '%#v'\n", i, c)
		if c.UnHashedKey != 0 {
			// Insert this element into new table
			cell, ok := s.Insert(c.UnHashedKey)
			if !ok {
				panic(fmt.Sprintf("key '%v' already exists in fresh table s: should be impossible", c.UnHashedKey))
			}
			*cell = *c
		}
	}

	VPrintf("\n ---- Done with Repopulate, now s = \n")
	VDump(s)

	t.DestroyHashTable()

	*t = *s
}

//----------------------------------------------
//  Iterator
//
//  sample use: given a HashTable h, enumerate h's contents with:
//
//		for it := offheap.NewIterator(h); it.Cur != nil; it.Next() {
//			found = append(found, it.Cur.UnHashedKey)
//		}
//
//----------------------------------------------

type Iterator struct {
	Tab *HashTable `capid:"0"`
	Pos int64      `capid:"1"`
	Cur *Cell      `capid:"2"` // will be set to nil when done with iteration.
}

func NewIterator(tab *HashTable) *Iterator {
	it := &Iterator{
		Tab: tab,
		Cur: &tab.ZeroCell,
		Pos: -1, // means we are at the ZeroCell to start with
	}

	if it.Tab.Population == 0 {
		it.Cur = nil
		it.Pos = -2
		return it
	}

	if !it.Tab.ZeroUsed {
		it.Next()
	}

	return it
}

func (it *Iterator) Done() bool {
	if it.Cur == nil {
		return true
	}
	return false
}

func (it *Iterator) Next() *Cell {

	// Already finished?
	if it.Cur == nil {
		return nil
	}

	// Iterate through the regular Cells
	it.Pos++
	for uint64(it.Pos) != it.Tab.ArraySize {
		it.Cur = it.Tab.CellAt(uint64(it.Pos))
		if it.Cur.UnHashedKey != 0 {
			return it.Cur
		}
		it.Pos++
	}

	// Finished
	it.Cur = nil
	it.Pos = -2
	return nil
}

func (t *HashTable) Dump() {
	for i := uint64(0); i < t.ArraySize; i++ {
		cell := t.CellAt(i)
		fmt.Printf("dump cell %d: \n cell.UnHashedKey: '%v'\n cell.ByteKey: '%s'\n cell.Value: '%#v'\n ===============", i, cell.UnHashedKey, string(cell.ByteKey[:]), cell.Value)
	}
}
