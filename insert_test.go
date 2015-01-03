package offheap_test

import (
	"testing"

	"github.com/glycerine/go-offheap-hashtable"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestInsert(t *testing.T) {

	h := offheap.NewHashTable(8)

	cv.Convey("inserting a non-zero key should enable retrieving them with Lookup", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)
		cv.So(h.Lookup(23), cv.ShouldEqual, nil)
		c, ok := h.Insert(23)
		c.Value = 55
		cv.So(c, cv.ShouldNotEqual, nil)
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)
		cv.So(h.Lookup(23), cv.ShouldNotEqual, nil)
		cell := h.Lookup(23)
		cv.So(cell.Value, cv.ShouldEqual, 55)
	})

	h.Clear()
	cv.Convey("inserting a zero key should also be retrievable with Lookup", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)
		cv.So(h.Lookup(0), cv.ShouldEqual, nil)
		c, ok := h.Insert(0)
		c.Value = 55
		cv.So(c, cv.ShouldNotEqual, nil)
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)
		cv.So(h.Lookup(0), cv.ShouldNotEqual, nil)
		cell := h.Lookup(0)
		cv.So(cell.Value, cv.ShouldEqual, 55)
	})

}
