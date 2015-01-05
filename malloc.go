package offheap

import (
	"fmt"
	"os"
	"syscall"

	"launchpad.net/gommap"
)

type MmapMalloc struct {
	Path         string
	File         *os.File
	Fd           int
	FileBytesLen int64
	BytesAlloc   int64
	MMap         gommap.MMap // equiv to Mem
	Mem          []byte      // equiv to Mmap
}

const MAP_ANONYMOUS = 0x20
const MAP_PRIVATE = 0x2
const MAP_SHARED = 0x1

const PROT_READ = 0x1
const PROT_WRITE = 0x2

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
func (mm *MmapMalloc) Free() {
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

		if DirExists(mm.Path) {
			panic(fmt.Sprintf("path '%s' already exists as a directory, so cannot be used as a memory mapped file.", mm.Path))
		}

		if !FileExists(mm.Path) {
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

	fmt.Printf("\n ------->> path = '%v',  mm.Fd = %v, with flags = %x, sz = %v,  prot = '%v'\n", path, mm.Fd, flags, sz, prot)

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

		var errno syscall.Errno
		if err != nil {
			errno = err.(syscall.Errno)
		}
		fmt.Printf("\n  mmap is '%#v', errno is '%v'\n", mmap, errno)
		panic(err)
	}

	mm.MMap = mmap
	mm.Mem = mmap

	return &mm
}

// returns once sync-to-disk is done
func (mm *MmapMalloc) BlockUntilSync() {
	mm.MMap.Sync(gommap.MS_SYNC)
}

// schedules sync, but may return before it is done.
func (mm *MmapMalloc) BackgrounSync() {
	mm.MMap.Sync(gommap.MS_ASYNC)
}
