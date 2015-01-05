package offheap_test

import (
	"testing"

	"github.com/glycerine/go-offheap-hashtable"
)

func TestMalloc(t *testing.T) {

	mm := offheap.Malloc(10*1024, "./mmap.dat")

	// regular bytes searching functions work, find the first line:
	//end := bytes.Index(mmap, []byte("\n"))
	//println(string([]byte(mmap[:end])))

	//	a := mmap[(end - 1):end]
	//	a[0]++ // change is saved to disk and read next time.

	writeme := []byte("hello memory mapped WORLD\n")
	copy(mm.Mem[0:26], writeme)
	mm.TruncateTo(int64(len(writeme)))

	mm.Free()

	mm2 := offheap.Malloc(10*1024, "")
	copy(mm2.Mem[0:26], writeme)
	mm2.Free()

}
