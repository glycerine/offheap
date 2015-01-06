package offheap

import (
	"fmt"
	"os"
	"syscall"

	"github.com/glycerine/gommap"
)

// provide Malloc() and Free() calls which request memory directly
// from the kernel via mmap(). Memory can optionally be backed
// by a file for simplicity/efficiency of saving to disk.
//
// For use when the Go GC overhead is too large, and you need to move
// the hash table off-heap.

type MmapMalloc struct {
	Path         string `capid:"0"`
	File         *os.File
	Fd           int         `capid:"1"`
	FileBytesLen int64       `capid:"2"`
	BytesAlloc   int64       `capid:"3"`
	MMap         gommap.MMap // equiv to Mem, just avoids casts everywhere.
	Mem          []byte      `capid:"4"` // equiv to Mmap
}

// only impacts the file underlying the mapping, not
// the mapping itself at this point.
func (mm *MmapMalloc) TruncateTo(newSize int64) {
	if mm.File == nil {
		panic("cannot call TruncateTo() on a non-file backed MmapMalloc.")
	}
	err := syscall.Ftruncate(int(mm.File.Fd()), newSize)
	if err != nil {
		panic(err)
	}
}

//
// offheap.Free()
//
// warning: any pointers still remaining will crash the program if dereferenced.
//
func (mm *MmapMalloc) Free() {
	if mm.File != nil {
		mm.File.Close()
	}
	err := mm.MMap.UnsafeUnmap()
	if err != nil {
		panic(err)
	}
}

//
// offheap.Malloc()
//
// If path is not empty then we memory map to the give path as well.
// Otherwise it is just like a call to malloc(): an anonymous memory allocation, outside the realm
// of the Go Garbage Collector.
//
// if numBytes is -1, then we take the size from the path file's size.
//
// return value:
//    .Mem holds a []byte pointing to the returned memory (as does .MMap, for use in other gommap calls).
//
func Malloc(numBytes int64, path string) *MmapMalloc {

	mm := MmapMalloc{
		Path: path,
	}

	flags := syscall.MAP_SHARED
	if path == "" {
		flags = syscall.MAP_ANON | syscall.MAP_PRIVATE
		mm.Fd = -1

		if numBytes < 0 {
			panic("numBytes was negative but path was also empty: don't know how much to allocate!")
		}

	} else {

		if dirExists(mm.Path) {
			panic(fmt.Sprintf("path '%s' already exists as a directory, so cannot be used as a memory mapped file.", mm.Path))
		}

		if !fileExists(mm.Path) {
			file, err := os.Create(mm.Path)
			if err != nil {
				panic(err)
			}
			mm.File = file
		} else {
			file, err := os.OpenFile(mm.Path, os.O_RDWR, 0777)
			if err != nil {
				panic(err)
			}
			mm.File = file
		}
		mm.Fd = int(mm.File.Fd())
	}

	sz := numBytes
	if path != "" {
		// file-backed memory
		if numBytes < 0 {

			var stat syscall.Stat_t
			if err := syscall.Fstat(mm.Fd, &stat); err != nil {
				panic(err)
			}
			sz = stat.Size

		} else {
			// set to the size requested
			err := syscall.Ftruncate(mm.Fd, numBytes)
			if err != nil {
				panic(err)
			}
		}
		mm.FileBytesLen = sz
	}
	// INVAR: sz holds non-negative length of memory/file.

	mm.BytesAlloc = sz

	prot := syscall.PROT_READ | syscall.PROT_WRITE

	VPrintf("\n ------->> path = '%v',  mm.Fd = %v, with flags = %x, sz = %v,  prot = '%v'\n", path, mm.Fd, flags, sz, prot)

	var mmap []byte
	var err error
	if mm.Fd == -1 {

		flags = syscall.MAP_ANON | syscall.MAP_PRIVATE
		mmap, err = syscall.Mmap(-1, 0, int(sz), prot, flags)

	} else {

		flags = syscall.MAP_SHARED
		mmap, err = syscall.Mmap(mm.Fd, 0, int(sz), prot, flags)
	}
	if err != nil {
		panic(err)
	}

	// duplicate member to avoid casts all over the place.
	mm.MMap = mmap
	mm.Mem = mmap

	return &mm
}

// BlockUntilSync() returns only once the file is synced to disk.
func (mm *MmapMalloc) BlockUntilSync() {
	mm.MMap.Sync(gommap.MS_SYNC)
}

// BackgroundSync() schedules a sync to disk, but may return before it is done.
// Without a call to either BackgroundSync() or BackgroundSync(), there
// is no guarantee that file has ever been written to disk at any point before
// the munmap() call that happens during Free(). See the man pages msync(2)
// and mmap(2) for details.
func (mm *MmapMalloc) BackgrounSync() {
	mm.MMap.Sync(gommap.MS_ASYNC)
}
