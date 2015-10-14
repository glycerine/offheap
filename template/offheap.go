package template

import (
	"fmt"
	"unsafe"

	"github.com/cheekybits/genny/generic"
	"github.com/remerge/offheap/util"
)

type _LT_ generic.Type
type _T_ generic.Type
type _C_ generic.Type

type EmptyStruct_LT_ struct{}

type HashTableMetadata_LT_ struct {
	MagicNumber uint64
	ArraySize   uint64
	Population  uint64
}

type HashTableCustomMetadata_LT_ struct {
	_C_
}

type HashTable_LT_ struct {
	*HashTableMetadata_LT_
	*HashTableCustomMetadata_LT_

	zeroCell uintptr
	cells    uintptr
	cellSize uintptr

	offheap      []byte
	offheapCells []byte
	mmm          util.MmapMalloc
}

// Create a new hash table, able to hold initialSize count of keys.
func NewHashTable_LT_(initialSize uint64) *HashTable_LT_ {
	return NewHashTable_LT_FileBacked(initialSize, "")
}

func NewHashTable_LT_FileBacked(initialSize uint64, filepath string) *HashTable_LT_ {
	metaSize := unsafe.Sizeof(HashTableMetadata_LT_{})
	cellSize := unsafe.Sizeof(Cell_LT_{})
	customSize := unsafe.Sizeof(HashTableCustomMetadata_LT_{})
	mmm := *util.Malloc(int64(metaSize+customSize+uintptr(initialSize+1)*cellSize), filepath)

	baseP := unsafe.Pointer(&mmm.Mem[0])
	base := (uintptr)(baseP)

	zeroCell := base + metaSize + customSize

	h := &HashTable_LT_{
		cellSize:     cellSize,
		offheap:      mmm.Mem,
		offheapCells: mmm.Mem[metaSize:],
		mmm:          mmm,
		zeroCell:     zeroCell,
		cells:        zeroCell + cellSize,
	}

	// check metadata
	h.HashTableMetadata_LT_ = (*HashTableMetadata_LT_)(baseP)
	if h.MagicNumber == 0x123456789ABCDEF {
		// mapped from file
	} else {
		// fresh
		h.MagicNumber = 0x123456789ABCDEF
		h.ArraySize = initialSize
		h.Population = 0
	}
	h.HashTableCustomMetadata_LT_ = (*HashTableCustomMetadata_LT_)((unsafe.Pointer)(base + metaSize))

	return h
}

type Cell_LT_ struct {
	unHashedKey uint64
	Value       _T_
}

var _empty_T_ _T_

// ZeroValue sets the cell's value to all zeros.
func (cell *Cell_LT_) ZeroValue() {
	*(&cell.Value) = _empty_T_
}

// Save syncs the memory mapped file to disk using MmapMalloc::BlockUntilSync()
func (t *HashTable_LT_) Save() {
	t.mmm.BlockUntilSync()
}

func (t *HashTable_LT_) cellAt(pos uint64) *Cell_LT_ {
	return (*Cell_LT_)(unsafe.Pointer(t.cells + uintptr(pos)*t.cellSize))
}

// DestroyHashTable frees the memory-mapping, returning the
// memory containing the hash table and its cells to the OS.
// By default the save-to-file-on-disk functionality in malloc.go is
// not used, but that can be easily activated. See malloc.go.
// Deferencing any cells/pointers into the hash table after
// destruction will result in crashing your process, almost surely.
func (t *HashTable_LT_) DestroyHashTable() {
	t.mmm.Free()
}

// Lookup a cell based on a uint64 key value. Returns nil if key not found.
func (t *HashTable_LT_) Lookup(key uint64) *Cell_LT_ {

	var cell *Cell_LT_

	if key == 0 {
		if t.zero().unHashedKey == 1 {
			return t.zero()
		}
		return nil
	}

	h := util.IntegerHash(uint64(key)) % t.ArraySize

	for {
		cell = t.cellAt(h)
		if cell.unHashedKey == key {
			return cell
		}
		if cell.unHashedKey == 0 {
			return nil
		}
		h++
		if h == t.ArraySize {
			h = 0
		}
	}
}

// Insert a key and get back the Cell_LT_ for that key, so
// as to enable assignment of Value within that Cell_LT_, for
// the specified key. The 2nd return value is false if
// key already existed (and thus required no addition); if
// the key already existed you can inspect the existing
// value in the *Cell_LT_ returned.

func (t *HashTable_LT_) zero() *Cell_LT_ {
	return (*Cell_LT_)(unsafe.Pointer(uintptr(t.zeroCell)))
}

func (t *HashTable_LT_) Insert(key uint64) (*Cell_LT_, bool) {
	var cell *Cell_LT_

	if key == 0 {
		zc := t.zero()
		// inuse
		isNew := zc.unHashedKey == 0
		if isNew {
			t.Population++
			zc.unHashedKey = 1
		}
		return zc, isNew
	}

	// if key != 0 {
	for {
		h := util.IntegerHash(uint64(key)) % t.ArraySize
		for {
			cell = t.cellAt(h)
			if cell.unHashedKey == key {
				// already exists
				return cell, false
			}
			if cell.unHashedKey == 0 {
				if (t.Population+1)*4 >= t.ArraySize*3 {
					t.Repopulate(t.ArraySize * 2)
					// resized, so start all over
					break
				}
				t.Population++
				cell.unHashedKey = key
				return cell, true
			}

			h++
			if h == t.ArraySize {
				h = 0
			}

		}
	}
}

// DeleteCell_LT_ deletes the cell pointed to by cell.
func (t *HashTable_LT_) DeleteCell_LT_(cell *Cell_LT_) {

	pos := uint64((uintptr(unsafe.Pointer(cell)) - uintptr(unsafe.Pointer(t.cells))) / uintptr(unsafe.Sizeof(Cell_LT_{})))

	// Delete from regular cells
	if pos < 0 || pos >= t.ArraySize {
		panic(fmt.Sprintf("cell out of bounds: pos %v was < 0 or >= t.ArraySize == %v", pos, t.ArraySize))
	}

	if t.cellAt(pos).unHashedKey == 0 {
		panic("zero unHashedKey in non-zero Cell_LT_!")
	}

	// Remove this cell by shuffling neighboring cells so there are no gaps in anyone's probe chain
	nei := pos + 1
	if nei >= t.ArraySize {
		nei = 0
	}
	var neighbor *Cell_LT_
	var circular_offset_ideal_pos int64
	var circular_offset_ideal_nei int64
	var cellPos *Cell_LT_

	for {
		neighbor = t.cellAt(nei)

		if neighbor.unHashedKey == 0 {
			// There's nobody to swap with. Go ahead and clear this cell, then return
			cellPos = t.cellAt(pos)
			cellPos.unHashedKey = 0
			cellPos.ZeroValue()
			t.Population--
			return
		}

		ideal := util.IntegerHash(neighbor.unHashedKey) % t.ArraySize

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
			*t.cellAt(pos) = *neighbor
			pos = nei
		}

		nei++
		if nei >= t.ArraySize {
			nei = 0
		}
	}
}

// Clear does not resize the table, but zeroes-out all entries.
func (t *HashTable_LT_) Clear() {
	// (Does not resize the array)
	// Clear regular cells

	for i := range t.offheapCells {
		t.offheapCells[i] = 0
	}
	t.Population = 0

	// Clear zero cell
	zc := t.zero()
	zc.unHashedKey = 0
	zc.ZeroValue()
}

// Compact will compress the hashtable so that it is at most
// 75% full.
func (t *HashTable_LT_) Compact() {
	t.Repopulate(util.UpperPowerOfTwo((t.Population*4 + 3) / 3))
}

// DeleteKey will delete the contents of the cell associated with key.
func (t *HashTable_LT_) DeleteKey(key uint64) {
	if key == 0 {
		zc := t.zero()
		if zc.unHashedKey == 1 {
			t.Population--
			zc.unHashedKey = 0
			zc.ZeroValue()
		}
		return
	}
	value := t.Lookup(key)
	if value != nil {
		t.DeleteCell_LT_(value)
	}
}

// Repopulate expands the hashtable to the desiredSize count of cells.
func (t *HashTable_LT_) Repopulate(desiredSize uint64) {
	if desiredSize&(desiredSize-1) != 0 {
		panic("desired size must be a power of 2")
	}
	if t.Population*4 > desiredSize*3 {
		panic("must have t.Population * 4  <= desiredSize * 3")
	}

	// Allocate new table
	s := NewHashTable_LT_(desiredSize)

	if t.zero().unHashedKey == 1 {
		zc := s.zero()
		zc.Value = t.zero().Value
		zc.unHashedKey = 1
		s.Population++
	}

	// Iterate through old table t, copy into new table s.
	var c *Cell_LT_

	for i := uint64(0); i < t.ArraySize; i++ {
		c = t.cellAt(i)
		if c.unHashedKey == 0 {
			continue
		}
		// Insert this element into new table
		cell, ok := s.Insert(c.unHashedKey)
		if !ok {
			panic(fmt.Sprintf("key@%d, '%v' already exists in fresh table s: should be impossible: %#v", i, c.unHashedKey, cell))
		}
		*cell = *c

	}
	t.DestroyHashTable()
	*t = *s
}

/*
_LT_Iterator

sample use: given a HashTable h, enumerate h's contents with:

    for it := offheap.NewIterator(h); it.Cur != nil; it.Next() {
      found = append(found, it.Cur.unHashedKey)
    }
*/
type _LT_Iterator struct {
	H   *HashTable_LT_
	Pos int64
	Cur *Cell_LT_
}

var _zeroCell_LT_ *Cell_LT_

// NewIterator creates a new iterator for HashTable tab.
func (h *HashTable_LT_) NewIterator() *_LT_Iterator {
	it := &_LT_Iterator{
		H:   h,
		Cur: &Cell_LT_{},
		Pos: -1, // means we are at the ZeroCell to start with
	}

	if it.H.Population == 0 {
		it.Cur = nil
		it.Pos = -2
		return it
	}

	zc := h.zero()
	if zc.unHashedKey == 1 {
		it.Cur = zc
	} else {
		it.Next()
	}
	return it
}

// Done checks to see if we have already iterated through all cells
// in the table. Equivalent to checking it.Cur == nil.
func (it *_LT_Iterator) Done() bool {
	if it.Cur == nil {
		return true
	}
	return false
}

// Next advances the iterator so that it.Cur points to the next
// filled cell in the table, and returns that cell. Returns nil
// once there are no more cells to be visited.
func (it *_LT_Iterator) Next() *Cell_LT_ {

	// Already finished?
	if it.Cur == nil {
		return nil
	}

	// Iterate through the regular cells
	it.Pos++
	for uint64(it.Pos) != it.H.ArraySize {
		it.Cur = it.H.cellAt(uint64(it.Pos))
		if it.Cur.unHashedKey != 0 {
			return it.Cur
		}
		it.Pos++
	}

	// Finished
	it.Cur = nil
	it.Pos = -2
	return nil
}
