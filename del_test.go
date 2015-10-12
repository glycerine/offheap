package offheap_test

import (
	"testing"

	"github.com/glycerine/offheap"
	cv "github.com/glycerine/goconvey/convey"
)

func TestDelete(t *testing.T) {

	h := offheap.NewHashTable(2)

	h.Clear()
	cv.Convey("delete should eliminate an element, leaving the rest", t, func() {
		N := 5

		// fill up a table
		cv.So(h.Population, cv.ShouldEqual, 0)
		for i := 0; i < N; i++ {
			a, ok := h.Insert(uint64(i))
			cv.So(ok, cv.ShouldEqual, true)
			cell := h.Lookup(uint64(i))
			cv.So(cell, cv.ShouldNotEqual, nil)
			cv.So(cell, cv.ShouldEqual, a)

			for j := i + 1; j < N+10; j++ {
				cell = h.Lookup(uint64(j))
				cv.So(cell, cv.ShouldEqual, nil)
			}
		}

		// now delete from it, checking all the way down for correctness
		for i := -10; i < N; i++ {
			h.DeleteKey(uint64(i))
			if i >= 0 {
				cell := h.Lookup(uint64(i))
				cv.So(cell, cv.ShouldEqual, nil)
				for j := i + 1; j < N; j++ {
					cell = h.Lookup(uint64(j))
					cv.So(cell, cv.ShouldNotEqual, nil)
					cv.So(cell.UnHashedKey, cv.ShouldEqual, j)
				}
			} else {
				cv.So(h.Population, cv.ShouldEqual, N)
			}
		}
		cell := h.Lookup(uint64(N + 1))
		cv.So(cell, cv.ShouldEqual, nil)

	})
}
