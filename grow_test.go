package offheap_test

import (
	"testing"

	"github.com/glycerine/offheap"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestGrowth(t *testing.T) {

	h := offheap.NewHashTable(2)

	cv.Convey("Given a size 2 table, inserting 0 and 1 should retain and recall the 0 and the 1 keys", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)
		for i := 0; i < 2; i++ {
			_, ok := h.Insert(uint64(i))
			cv.So(ok, cv.ShouldEqual, true)
		}
		cv.So(h.Population, cv.ShouldEqual, 2)
		cv.So(h.ArraySize, cv.ShouldEqual, 4)
	})

	h.Clear()
	cv.Convey("inserting more than the current size should automatically grow the table", t, func() {
		N := 100

		cv.So(h.Population, cv.ShouldEqual, 0)
		for i := 0; i < N; i++ {
			_, ok := h.Insert(uint64(i))
			cv.So(ok, cv.ShouldEqual, true)
		}
		cv.So(h.Population, cv.ShouldEqual, N)
		cv.So(h.ArraySize, cv.ShouldEqual, upper_power_of_two(uint64(float64(N)*4.0/3.0)))

		for i := 0; i < N; i++ {
			cell := h.Lookup(uint64(i))
			cv.So(cell, cv.ShouldNotEqual, nil)
			cv.So(cell.UnHashedKey, cv.ShouldEqual, i)
		}
		cell := h.Lookup(uint64(N + 1))
		cv.So(cell, cv.ShouldEqual, nil)

	})
}

func upper_power_of_two(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}
