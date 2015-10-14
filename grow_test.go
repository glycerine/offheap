package offheap

import (
	"testing"

	"github.com/remerge/offheap/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGrowth(t *testing.T) {

	h := NewHashTableSomeStruct(2)

	Convey("Given a size 2 table, inserting 0 and 1 should retain and recall the 0 and the 1 keys", t, func() {
		So(h.Population, ShouldEqual, 0)
		for i := 0; i < 2; i++ {
			_, ok := h.Insert(uint64(i))
			So(ok, ShouldEqual, true)
		}
		So(h.Population, ShouldEqual, 2)
		So(h.ArraySize, ShouldEqual, 4)
	})

	h.Clear()
	Convey("inserting more than the current size should automatically grow the table", t, func() {
		N := uint64(100)

		So(h.Population, ShouldEqual, 0)
		for i := uint64(0); i < N; i++ {
			_, ok := h.Insert(i)
			So(ok, ShouldEqual, true)
		}
		So(h.Population, ShouldEqual, N)
		So(h.ArraySize, ShouldEqual, util.UpperPowerOfTwo(uint64(float64(N)*4.0/3.0)))

		cell := h.Lookup(0)
		So(cell, ShouldNotEqual, nil)
		So(cell.unHashedKey, ShouldEqual, 1)

		for i := uint64(1); i < N; i++ {
			cell := h.Lookup(i)
			So(cell, ShouldNotEqual, nil)
			So(cell.unHashedKey, ShouldEqual, i)
		}
		cell = h.Lookup(N + 1)
		So(cell, ShouldEqual, nil)

	})
}
