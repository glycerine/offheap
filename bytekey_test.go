package offheap_test

import (
	"testing"

	"github.com/glycerine/go-offheap-hashtable"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestByteKeyInsertLookup(t *testing.T) {

	h := offheap.NewHashTable(8)

	cv.Convey("inserting with a byte slice key using InsertBK should enable retrieving them with LookupBK", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)

		look, ok := h.LookupBK([]byte("hello"))
		cv.So(ok, cv.ShouldEqual, false)
		cv.So(look, cv.ShouldEqual, interface{}(nil))

		ok = h.InsertBK([]byte("hello"), "world")
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)

		val, ok := h.LookupBK([]byte("hello"))
		cv.So(val.(string), cv.ShouldEqual, "world")
		cv.So(ok, cv.ShouldEqual, true)
	})

}

func TestStringKeyInsertLookup(t *testing.T) {

	h := offheap.NewHashTable(8)

	cv.Convey("inserting with a byte slice key using InsertBK should enable retrieving them with LookupBK", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)

		look, ok := h.LookupStringKey("hello")
		cv.So(ok, cv.ShouldEqual, false)
		cv.So(look, cv.ShouldEqual, interface{}(nil))

		ok = h.InsertStringKey("hello", "world")
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)

		h.Dump()

		val, ok := h.LookupStringKey("hello")
		cv.So(val.(string), cv.ShouldEqual, "world")
		cv.So(ok, cv.ShouldEqual, true)
	})

}
