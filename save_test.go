package offheap_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/glycerine/offheap"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestSaveRestore(t *testing.T) {

	fn := "save_test_hash.ohdat"
	os.Remove(fn)
	h := offheap.NewHashFileBacked(8, fn)

	cv.Convey("saving and then loading a table should restore the contents from disk", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)
		cv.So(h.Lookup(23), cv.ShouldEqual, nil)
		c, ok := h.Insert(23)
		c.SetInt(55)
		cv.So(c, cv.ShouldNotEqual, nil)
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)
		cv.So(h.Lookup(23), cv.ShouldNotEqual, nil)
		cell := h.Lookup(23)
		cv.So(cell.Value[0], cv.ShouldEqual, 55)

		h.InsertIntValue(45, 28)

		// h has:
		// 23 -> 55
		// 45 -> 28

		h.Save()

		// copy to a new file to be sure everything is there, then mmap the new file
		fncopy := fn + ".copy"
		err := exec.Command("/bin/cp", "-p", fn, fncopy).Run()
		if err != nil {
			panic(err)
		}

		h2 := offheap.NewHashFileBacked(8, fncopy)

		cell2 := h2.Lookup(23)
		cv.So(cell2.Value[0], cv.ShouldEqual, 55)

		cell2 = h2.Lookup(45)
		cv.So(cell2.Value[0], cv.ShouldEqual, 28)

	})
}
