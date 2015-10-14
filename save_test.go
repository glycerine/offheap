package offheap

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSaveRestore(t *testing.T) {

	fn := "save_test_hash.ohdat"
	err := os.Remove(fn)
	if err != nil && !strings.HasSuffix(err.Error(), "no such file or directory") {
		panic(err)
	}
	defer os.Remove(fn)

	h := NewHashTableIntFileBacked(8, fn)

	Convey("saving and then loading a table should restore the contents from disk", t, func() {
		So(h.Population, ShouldEqual, 0)
		So(h.Lookup(23), ShouldEqual, nil)
		h.InsertInt(23, 55)
		So(h.Population, ShouldEqual, 1)
		So(h.Lookup(23), ShouldNotEqual, nil)
		cell := h.Lookup(23)
		So(cell.Value, ShouldEqual, 55)

		h.InsertInt(45, 28)
		h.InsertInt(0, 111)

		// h has:
		// 23 -> 55
		// 45 -> 28
		// 0 -> 111

		// custom metadata
		h.A = 42
		h.B = 33.33
		copy(h.C[:], []byte("1234567890"))
		h.Save()

		// copy to a new file to be sure everything is there, then mmap the new file
		fncopy := fn + ".copy"
		err := os.Remove(fncopy)
		if err != nil && !strings.HasSuffix(err.Error(), "no such file or directory") {
			panic(err)
		}
		err = exec.Command("/bin/cp", "-p", fn, fncopy).Run()
		if err != nil {
			panic(err)
		}
		defer os.Remove(fncopy)

		h2 := NewHashTableIntFileBacked(8, fn)

		So(h2.Population, ShouldEqual, 3)

		v, _ := h2.LookupInt(23)
		So(v, ShouldEqual, 55)
		v, _ = h2.LookupInt(45)
		So(v, ShouldEqual, 28)
		v, _ = h2.LookupInt(0)
		So(v, ShouldEqual, 111)

		So(h2.A, ShouldEqual, 42)
		So(h2.B, ShouldEqual, 33.33)

		So(string(h2.C[:]), ShouldEqual, "1234567890")

	})
}
