package offheap_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/offheap"
)

func TestSaveRestore(t *testing.T) {

	fn := "save_test_hash.ohdat"
	err := os.Remove(fn)
	if err != nil && !strings.HasSuffix(err.Error(), "no such file or directory") {
		panic(err)
	}
	defer os.Remove(fn)

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

		cv.So(h.Population, cv.ShouldEqual, 2)
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

		h2 := offheap.NewHashFileBacked(8, fncopy)

		cell2 := h2.Lookup(23)
		cv.So(cell2.Value[0], cv.ShouldEqual, 55)

		cell2 = h2.Lookup(45)
		cv.So(cell2.Value[0], cv.ShouldEqual, 28)

		cv.So(h2.Population, cv.ShouldEqual, 2)
	})
}

func TestMetaSaveRestoreMetadataInMsgpack(t *testing.T) {
	cv.Convey("saving the metadata of a table using msgpack should result in restore-able metadata", t, func() {
		t1 := offheap.HashTable{
			Population: 4,
			ZeroUsed:   true,
			ZeroCell: offheap.Cell{
				UnHashedKey: 43,
				ByteKey:     offheap.Key_t{0x23},
				Value:       offheap.Val_t{0x57},
			},
		}
		bts, err := t1.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}

		t2 := offheap.HashTable{}
		left, err := t2.UnmarshalMsg(bts)
		if err != nil {
			t.Fatal(err)
		}
		if len(left) > 0 {
			t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
		}

		cv.So(t2, cv.ShouldResemble, t1)
		fmt.Printf("\n len(bts) = %d\n", len(bts))

		// if we ever get bigger than this, there will be problems as the
		// metadata will overflow into the Cell data area on serialization.
		// This is not a perfect
		// test since msgpack serialization sizes vary depending on content. Hence we
		// give our selves an extra 256 byte buffer.
		cv.So(len(bts), cv.ShouldBeLessThan, offheap.MetadataHeaderMaxBytes-256)
	})
}
