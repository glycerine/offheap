package offheap

import (
	"fmt"
	"unsafe"
)

// Copyright (C) 2015 by Jason E. Aten, Ph.D.
//
// Inspired by the public domain C++ code of
//    https://github.com/preshing/CompareIntegerMaps
// See also
//    http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/
// for performance studies.
//

//----------------------------------------------
//  HashTable
//
//  Maps pointer-sized integers to pointer-sized integers.
//  Uses open addressing with linear probing.
//  In the t.cells array, HashedKey = 0 is reserved to indicate an unused cell.
//  Actual value for key 0 (if any) is stored in t.zeroCell.
//  The hash table automatically doubles in size when it becomes 75% full.
//  The hash table never shrinks in size, even after Clear(), unless you explicitly
//  call Compact().
//----------------------------------------------

type Cell struct {
	HashedKey uint64
	Value     interface{}
}

type HashTable struct {
	Cells      []Cell
	ArraySize  uint64
	Population uint64
	ZeroUsed   bool
	ZeroCell   Cell
}

func NewHashTable(initialSize uint64) *HashTable {
	return &HashTable{
		// todo: allocate this off-heap instead
		Cells:     make([]Cell, initialSize),
		ArraySize: initialSize,
	}
}

func (t *HashTable) DestroyHashTable() {
	// todo: release the off-heap allocation here
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
			cell = &(t.Cells[h])
			if cell.HashedKey == key {
				return cell
			}
			if cell.HashedKey == 0 {
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
				cell = &(t.Cells[h])

				if cell.HashedKey == key {
					// already exists
					return cell, false
				}
				if cell.HashedKey == 0 {
					if (t.Population+1)*4 >= t.ArraySize*3 {
						VPrintf("detected (t.Population+1)*4 >= t.ArraySize*3, i.e. %v >= %v, calling Repop with double the size\n", (t.Population+1)*4, t.ArraySize*3)
						t.Repopulate(t.ArraySize * 2)
						// resized, so start all over
						break
					}
					t.Population++
					cell.HashedKey = key
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

func (t *HashTable) InsertIntValue(key uint64, value int) bool {
	cell, ok := t.Insert(key)
	cell.Value = value
	return ok
}

func (t *HashTable) DeleteCell(cell *Cell) {

	if cell == &t.ZeroCell {
		// Delete zero cell
		if !t.ZeroUsed {
			panic("deleting zero element when not used")
		}
		t.ZeroUsed = false
		cell.Value = nil
		t.Population--
		return

	} else {

		pos := uint64((uintptr(unsafe.Pointer(cell)) - uintptr(unsafe.Pointer(&t.Cells[0]))) / uintptr(unsafe.Sizeof(Cell{})))

		// Delete from regular Cells
		if pos < 0 || pos >= t.ArraySize {
			panic(fmt.Sprintf("cell out of bounds: pos %v was < 0 or >= t.ArraySize == %v", pos, t.ArraySize))
		}
		if t.Cells[pos].HashedKey == 0 {
			panic("zero HashedKey in non-zero Cell!")
		}

		// Remove this cell by shuffling neighboring Cells so there are no gaps in anyone's probe chain
		nei := pos + 1
		if nei >= t.ArraySize {
			nei = 0
		}
		var neighbor *Cell
		var circular_offset_ideal_pos int64
		var circular_offset_ideal_nei int64

		for {
			neighbor = &t.Cells[nei]

			if neighbor.HashedKey == 0 {
				// There's nobody to swap with. Go ahead and clear this cell, then return
				t.Cells[pos].HashedKey = 0
				t.Cells[pos].Value = nil
				t.Population--
				return
			}

			ideal := integerHash(neighbor.HashedKey) % t.ArraySize

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
				t.Cells[pos] = *neighbor
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

	// todo, change to use off heap memory
	for i := range t.Cells {
		t.Cells[i] = Cell{}
	}
	t.Population = 0

	// Clear zero cell
	t.ZeroUsed = false
	t.ZeroCell.Value = 0
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

	// Get start/end pointers of old array
	oldCells := t.Cells

	// Allocate new array
	t.ArraySize = desiredSize
	t.Cells = make([]Cell, t.ArraySize)

	// Iterate through old array
	// (any zero entry can stay in place; so ignore HashedKey == 0 below).
	var c *Cell
	var pos uint64
	for i := range oldCells {
		{
			c = &oldCells[i]
			VPrintf("\n in oldCell copy loop, at i = %v, and c = '%#v'\n", i, c)
			if c.HashedKey != 0 {
				// Insert this element into new array
				pos = integerHash(c.HashedKey) % t.ArraySize

				VPrintf("   in Repop, pos = %v for c.HashedKey = %v and t.ArraySize = %v\n", pos, c.HashedKey, t.ArraySize)
				for {
					cell := &t.Cells[pos]
					VPrintf("cell = %v, pos = %v, t.Cells = %v\n", cell, pos, t.Cells)

					if cell.HashedKey == 0 {
						// Insert here
						*cell = *c
						break
					}
					pos++
					if pos >= t.ArraySize {
						pos = 0
					}
				}
			}
		}

		// Delete old array; happens when oldCells goes out of scope
		// todo: delete in off-heap space
	}

	VPrintf("\n ---- Done with Repopulate, now t = \n")
	VDump(t)
}

//----------------------------------------------
//  Iterator
//----------------------------------------------

type Iterator struct {
	Tab *HashTable
	Pos int64
	Cur *Cell // nil when done
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

func (it *Iterator) Next() *Cell {

	// Already finished?
	if it.Cur == nil {
		return nil
	}

	// Iterate through the regular Cells
	it.Pos++
	for uint64(it.Pos) != it.Tab.ArraySize {
		it.Cur = &it.Tab.Cells[it.Pos]
		if it.Cur.HashedKey != 0 {
			return it.Cur
		}
		it.Pos++
	}

	// Finished
	it.Cur = nil
	it.Pos = -2
	return nil
}
