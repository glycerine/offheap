package offheap

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

//go:generate msgp

// Copyright (C) 2015 by Jason E. Aten, Ph.D.

// metadata serialization size can never grow bigger
// than MetadataHeaderMaxBytes (currently one 4K page), without impacting
// backwards compatibility. We reserve this many bytes
// at the beginning of the memory mapped file for the metadata.
const MetadataHeaderMaxBytes = 4096

// HashTable represents the off-heap hash table.
// Create a new one with NewHashTable(), then use
// Lookup(), Insert(), and DeleteKey() on it.
// HashTable is meant to be customized by the
// user, to reflect your choice of key and value
// types. See StringHashTable and ByteKeyHashTable
// for examples of this specialization process.
type HashTable struct {
	MagicNumber   int     // distinguish between reading from empty file versus an on-disk HashTable.
	Cells         uintptr `msg:"-"`
	CellSz        uint64
	ArraySize     uint64
	Population    uint64
	ZeroUsed      bool
	ZeroCell      Cell
	OffheapHeader []byte     `msg:"-"`
	OffheapCells  []byte     `msg:"-"`
	Mmm           MmapMalloc `msg:"-"`
}

// Create a new hash table, able to hold initialSize count of keys.
func NewHashTable(initialSize int64) *HashTable {
	return NewHashFileBacked(initialSize, "")
}

func NewHashFileBacked(initialSize int64, filepath string) *HashTable {
	//fmt.Printf("\n\n  NewHashFileBacked called!\n\n")
	t := HashTable{
		MagicNumber: 3030675468910466832,
		CellSz:      uint64(unsafe.Sizeof(Cell{})),
	}

	// off-heap and off-gc version
	if initialSize < 0 {
		t.Mmm = *Malloc(-1, filepath)
	} else {
		t.ArraySize = uint64(initialSize)
		t.Mmm = *Malloc(t.bytesFromArraySizeAndHeader(t.ArraySize), filepath)
	}
	//p("&t.Mmm.Mem[0] = %p\n", &t.Mmm.Mem[0])
	t.OffheapHeader = t.Mmm.Mem
	t.OffheapCells = t.Mmm.Mem[MetadataHeaderMaxBytes:]
	t.Cells = (uintptr)(unsafe.Pointer(&t.OffheapCells[0]))

	// test deserialize
	t2 := HashTable{}
	_, err := t2.UnmarshalMsg(t.OffheapHeader)
	if err != nil {
		//fmt.Printf("UnmarshalMsg err = %v\n", err)
		// common to not be able to deserialize 0 bytes, don't worry
	} else {
		//fmt.Printf("\n deserialized okay! t2=%#v\n t=%#v\n", t2, t)
		// no error during de-serialize; okay but verify our MagicNumer too.
		if t2.MagicNumber == t.MagicNumber {
			// okay to copy in, we are restoring an existing table from disk
			t.Population = t2.Population
			t.ZeroUsed = t2.ZeroUsed
			t.ZeroCell = t2.ZeroCell
		}
	}
	// old: off-gc but still on-heap version
	//	t.ArraySize = initialSize
	//	t.Offheap = make([]byte, t.ArraySize*t.CellSz)
	//	t.Cells = (uintptr)(unsafe.Pointer(&t.Offheap[0]))

	// on-heap version:
	//Cells:     make([]Cell, initialSize),

	return &t
}

// Key_t is the basic type for keys. Users of the library will
// probably redefine this.
type Key_t [64]byte

// Val_t is the basic type for values stored in the cells in the table.
// Users of the library will probably redefine this to be a different
// size at the very least.
type Val_t [56]byte

// Cell is the basic payload struct, stored inline in the HashTable. The
// cell is returned by the fundamental Lookup() function. The member
// Value is where the value that corresponds to the key (in ByteKey)
// is stored. Both the key (in ByteKey) and the value (in Value) are
// stored inline inside the hashtable, so that all storage for
// the hashtable is in the same offheap segment. The uint64 key given
// to fundamental Insert() method is stored in UnHashedKey. The hashed value of the
// UnHashedKey is not stored in the Cell, but rather computed as needed
// by the basic Insert() and Lookup() methods.
type Cell struct {
	UnHashedKey uint64
	ByteKey     Key_t
	Value       Val_t // customize this to hold your value's data type entirely here.
}

/*
SetValue stores any value v in the Cell. Note that
users of the library will need to extend this for
their type. Only strings of length less than 56,
and integers are handled by default.
*/
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

// ZeroValue sets the cell's value to all zeros.
func (cell *Cell) ZeroValue() {
	for i := range cell.Value[:] {
		cell.Value[i] = 0
	}
}

// SetString stores string s (up to val_t length, currently 56 bytes) in cell.Value.
func (cell *Cell) SetString(s string) {
	copy(cell.Value[:], []byte(s))
}

// GetString retreives a string value from the cell.Value.
func (cell *Cell) GetString() string {
	return string([]byte(cell.Value[:]))
}

// SetInt stores an integer value in the cell.
func (cell *Cell) SetInt(n int) {
	binary.LittleEndian.PutUint64(cell.Value[:8], uint64(n))
}

// GetInt retreives an integer value from the cell.
func (cell *Cell) GetInt() int {
	return int(binary.LittleEndian.Uint64(cell.Value[:8]))
}

// SetInt sets an int value for Val_t v.
func (v *Val_t) SetInt(n int) {
	binary.LittleEndian.PutUint64((*v)[:8], uint64(n))
}

// GetInt gets an int value for Val_t v.
func (v *Val_t) GetInt() int {
	return int(binary.LittleEndian.Uint64((*v)[:8]))
}

// SetString sets a string value for Val_t v.
func (v *Val_t) SetString(s string) {
	copy((*v)[:], []byte(s))
}

// GetString retreives a string value for Val_t v.
func (v *Val_t) GetString() string {
	return string([]byte((*v)[:]))
}

// Save syncs the memory mapped file to disk using MmapMalloc::BlockUntilSync().
// If background is true, the save using BackgroundSync() instead of blocking.
func (t *HashTable) Save(background bool) error {
	bts, err := t.MarshalMsg(nil)
	if err != nil {
		return fmt.Errorf("offheap.HashTable.Save() error: serialization error from msgp.MarshalMsg(): '%s'", err)
	}

	if len(bts) > MetadataHeaderMaxBytes {
		return fmt.Errorf("offheap.HashTable.Save() error: serialization too long: len(bts)==%d is > MetadataHeaderMaxBytes==%d, so serializing the HashTable metadata would overwrite the cell contents; rather than corrupting the cell data we are returning an error here.", len(bts), MetadataHeaderMaxBytes)
	}
	nw := copy(t.OffheapHeader, bts)
	_ = nw
	//fmt.Printf("saved bts header of size %v by writing %v = '%s'. t.OffheapHeader = %p, t.OffheapCells = %p, &t.Mmm.Mem[0] = %p\n", len(bts), nw, string(bts), &t.OffheapHeader[0], &t.OffheapCells[0], &t.Mmm.Mem[0])

	if background {
		t.Mmm.BackgroundSync()
		//fmt.Printf("\ndone with call to initiate background sync.\n")
	} else {
		t.Mmm.BlockUntilSync()
		//fmt.Printf("\ndone with block until sync\n")
	}
	return nil
}

// CellAt: fetch the cell at a given index. E.g. t.CellAt(pos) replaces t.Cells[pos]
func (t *HashTable) CellAt(pos uint64) *Cell {

	// off heap version
	return (*Cell)(unsafe.Pointer(uintptr(t.Cells) + uintptr(pos*t.CellSz)))

	// on heap version, back when t.Cells was []Cell
	//return &(t.Cells[pos])
}

// DestroyHashTable frees the memory-mapping, returning the
// memory containing the hash table and its cells to the OS.
// By default the save-to-file-on-disk functionality in malloc.go is
// not used, but that can be easily activated. See malloc.go.
// Deferencing any cells/pointers into the hash table after
// destruction will result in crashing your process, almost surely.
func (t *HashTable) DestroyHashTable() {
	t.Mmm.Free()
}

// Lookup a cell based on a uint64 key value. Returns nil if key not found.
func (t *HashTable) Lookup(key uint64) *Cell {

	var cell *Cell

	if key == 0 {
		if t.ZeroUsed {
			return &t.ZeroCell
		}
		return nil

	} else {
		//p("for t = %p, t.ArraySize = %v", t, t.ArraySize)
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

// Insert a key and get back the Cell for that key, so
// as to enable assignment of Value within that Cell, for
// the specified key. The 2nd return value is false if
// key already existed (and thus required no addition); if
// the key already existed you can inspect the existing
// value in the *Cell returned.
func (t *HashTable) Insert(key uint64) (*Cell, bool) {

	vprintf("\n ---- Insert(%v) called with t = \n", key)
	vdump(t)

	defer func() {
		vprintf("\n ---- Insert(%v) done, with t = \n", key)
		vdump(t)
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
						vprintf("detected (t.Population+1)*4 >= t.ArraySize*3, i.e. %v >= %v, calling Repop with double the size\n", (t.Population+1)*4, t.ArraySize*3)
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

// InsertIntValue inserts value under key in the table.
func (t *HashTable) InsertIntValue(key uint64, value int) bool {
	cell, ok := t.Insert(key)
	cell.SetValue(value)
	return ok
}

// DeleteCell deletes the cell pointed to by cell.
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

// Clear does not resize the table, but zeroes-out all entries.
func (t *HashTable) Clear() {
	// (Does not resize the array)
	// Clear regular Cells

	for i := range t.OffheapCells {
		t.OffheapCells[i] = 0
	}
	t.Population = 0

	// Clear zero cell
	t.ZeroUsed = false
	t.ZeroCell.ZeroValue()
}

// Compact will compress the hashtable so that it is at most
// 75% full.
func (t *HashTable) Compact() {
	t.Repopulate(upper_power_of_two((t.Population*4 + 3) / 3))
}

// DeleteKey will delete the contents of the cell associated with key.
func (t *HashTable) DeleteKey(key uint64) {
	value := t.Lookup(key)
	if value != nil {
		t.DeleteCell(value)
	}
}

// Repopulate expands the hashtable to the desiredSize count of cells.
func (t *HashTable) Repopulate(desiredSize uint64) {

	//p("top of Repopulate(%v)", desiredSize)
	vprintf("\n ---- Repopulate called with t = \n")
	vdump(t)

	if desiredSize&(desiredSize-1) != 0 {
		panic("desired size must be a power of 2")
	}
	if t.Population*4 > desiredSize*3 {
		panic("must have t.Population * 4  <= desiredSize * 3")
	}

	// Allocate new table
	// TODO: implement growmap for mmap backed resizing.
	s := NewHashTable(int64(desiredSize))
	s.ZeroUsed = t.ZeroUsed
	if t.ZeroUsed {
		s.ZeroCell = t.ZeroCell
		s.Population++
	}

	// Iterate through old table t, copy into new table s.
	var c *Cell

	for i := uint64(0); i < t.ArraySize; i++ {
		c = t.CellAt(i)
		vprintf("\n in oldCell copy loop, at i = %v, and c = '%#v'\n", i, c)
		if c.UnHashedKey != 0 {
			// Insert this element into new table
			cell, ok := s.Insert(c.UnHashedKey)
			if !ok {
				panic(fmt.Sprintf("key '%v' already exists in fresh table s: should be impossible", c.UnHashedKey))
			}
			*cell = *c
		}
	}

	vprintf("\n ---- Done with Repopulate, now s = \n")
	vdump(s)

	t.DestroyHashTable()

	*t = *s
	//p("bottom of Repopulate(), now t.ArraySize = %v, t = %p", t.ArraySize, t)
}

/*
Iterator

sample use: given a HashTable h, enumerate h's contents with:

    for it := offheap.NewIterator(h); it.Cur != nil; it.Next() {
      found = append(found, it.Cur.UnHashedKey)
    }
*/
type Iterator struct {
	Tab *HashTable
	Pos int64
	Cur *Cell // will be set to nil when done with iteration.
}

// NewIterator creates a new iterator for HashTable tab.
func (tab *HashTable) NewIterator() *Iterator {
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

// Done checks to see if we have already iterated through all cells
// in the table. Equivalent to checking it.Cur == nil.
func (it *Iterator) Done() bool {
	if it.Cur == nil {
		return true
	}
	return false
}

// Next advances the iterator so that it.Cur points to the next
// filled cell in the table, and returns that cell. Returns nil
// once there are no more cells to be visited.
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

// Dump provides a diagnostic dump of the full HashTable contents.
func (t *HashTable) Dump() {
	for i := uint64(0); i < t.ArraySize; i++ {
		cell := t.CellAt(i)
		fmt.Printf("dump cell %d: \n cell.UnHashedKey: '%v'\n cell.ByteKey: '%s'\n cell.Value: '%#v'\n ===============", i, cell.UnHashedKey, string(cell.ByteKey[:]), cell.Value)
	}
}

// InsertBK is the insert function for []byte keys.
// By default only len(Key_t) bytes are used in the key.
func (t *HashTable) InsertBK(bytekey []byte, value interface{}) bool {
	return ((*ByteKeyHashTable)(t)).InsertBK(bytekey, value)
}

func (t *HashTable) LookupBK(bytekey []byte) (Val_t, bool) {
	return ((*ByteKeyHashTable)(t)).LookupBK(bytekey)
}

func (t *HashTable) LookupBKInt(bytekey []byte) (int, bool) {
	v, found := ((*ByteKeyHashTable)(t)).LookupBK(bytekey)
	if !found {
		return -1, found
	}
	return v.GetInt(), true
}

func (t *HashTable) bytesFromArraySizeAndHeader(arraySize uint64) int64 {
	return int64(arraySize*t.CellSz) + MetadataHeaderMaxBytes
}
