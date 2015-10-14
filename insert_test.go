package offheap

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type SomeStruct struct {
	A int
	B int
}

func TestInsert(t *testing.T) {

	h := NewHashTableSomeStruct(8)

	Convey("inserting a non-zero key should enable retrieving them with Lookup", t, func() {
		So(h.Population, ShouldEqual, 0)
		So(h.Lookup(23), ShouldEqual, nil)
		c, ok := h.Insert(23)
		c.Value.A = 5
		c.Value.B = 10
		So(c, ShouldNotEqual, nil)
		So(ok, ShouldEqual, true)
		So(h.Population, ShouldEqual, 1)
		So(h.Lookup(23), ShouldNotEqual, nil)
		cell := h.Lookup(23)
		So(cell.Value.A, ShouldEqual, 5)
		So(cell.Value.B, ShouldEqual, 10)
	})

	h.Clear()
	Convey("inserting a zero key should also be retrievable with Lookup", t, func() {
		So(h.Population, ShouldEqual, 0)
		So(h.Lookup(0), ShouldEqual, nil)
		c, ok := h.Insert(0)
		c.Value.A = 5
		So(c, ShouldNotEqual, nil)
		So(ok, ShouldEqual, true)
		So(h.Population, ShouldEqual, 1)
		So(h.Lookup(0), ShouldNotEqual, nil)
		cell := h.Lookup(0)
		So(cell.Value.A, ShouldEqual, 5)
	})

	h.Clear()
	Convey("Insert()-ing the same key twice should return false for the 2nd param on encountering the same key again. For 0 key.", t, func() {
		So(h.Population, ShouldEqual, 0)
		c, ok := h.Insert(0)
		So(c, ShouldNotEqual, nil)
		So(c.unHashedKey, ShouldEqual, 1)
		So(ok, ShouldEqual, true)

		c, ok = h.Insert(0)
		So(c, ShouldNotEqual, nil)
		So(c.unHashedKey, ShouldEqual, 1)
		So(ok, ShouldEqual, false)
	})

	h.Clear()
	Convey("Insert()-ing the same key twice should return false for the 2nd param on encountering the same key again. For not-zero key.", t, func() {
		So(h.Population, ShouldEqual, 0)
		c, ok := h.Insert(1)
		So(c, ShouldNotEqual, nil)
		So(c.unHashedKey, ShouldEqual, 1)
		So(ok, ShouldEqual, true)

		c, ok = h.Insert(1)
		So(c, ShouldNotEqual, nil)
		So(c.unHashedKey, ShouldEqual, 1)
		So(ok, ShouldEqual, false)
	})

}
