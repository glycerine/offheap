package offheap

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDelete(t *testing.T) {

	h := NewHashTableInt(2)

	h.Clear()
	Convey("delete should eliminate an element, leaving the rest", t, func() {
		N := uint64(5)

		// fill up a table
		So(h.Population, ShouldEqual, 0)
		for i := uint64(0); i < N; i++ {
			a, ok := h.Insert(i)
			So(ok, ShouldEqual, true)
			cell := h.Lookup(i)
			So(cell, ShouldNotEqual, nil)
			So(cell, ShouldEqual, a)

			for j := i + 1; j < N+10; j++ {
				cell = h.Lookup(j)
				So(cell, ShouldEqual, nil)
			}
		}

		// now delete from it, checking all the way down for correctness
		for i := uint64(0); i < N; i++ {
			h.DeleteKey(i)
			if i >= 0 {
				cell := h.Lookup(i)
				So(cell, ShouldEqual, nil)
				for j := i + 1; j < N; j++ {
					cell = h.Lookup(uint64(j))
					So(cell, ShouldNotEqual, nil)
					So(cell.unHashedKey, ShouldEqual, j)
				}
			} else {
				So(h.Population, ShouldEqual, N)
			}
		}
		cell := h.Lookup(uint64(N + 1))
		So(cell, ShouldEqual, nil)

	})
}
