package offheap_test

import (
	"testing"

	"github.com/glycerine/offheap"
	cv "github.com/glycerine/goconvey/convey"
)

func TestByteKeyInsertLookup(t *testing.T) {

	h := offheap.NewByteKeyHashTable(8)

	cv.Convey("inserting with a byte slice key using InsertBK should enable retrieving them with LookupBK", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)

		look, ok := h.LookupBK([]byte("hello"))
		cv.So(ok, cv.ShouldEqual, false)
		cv.So(look, cv.ShouldResemble, offheap.Val_t{})

		ok = h.InsertBK([]byte("hello"), "world")
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)

		val, ok := h.LookupBK([]byte("hello"))
		cv.So(val.GetString(), cv.ShouldStartWith, "world")
		cv.So(ok, cv.ShouldEqual, true)
	})

}

func TestStringKeyInsertLookup(t *testing.T) {

	h := offheap.NewStringHashTable(8)

	cv.Convey("inserting with a byte slice key using InsertBK should enable retrieving them with LookupBK", t, func() {
		cv.So(h.Population, cv.ShouldEqual, 0)

		look, ok := h.LookupStringKey("hello")
		cv.So(ok, cv.ShouldEqual, false)
		cv.So(look, cv.ShouldResemble, offheap.Val_t{})

		ok = h.InsertStringKey("hello", "world")
		cv.So(ok, cv.ShouldEqual, true)
		cv.So(h.Population, cv.ShouldEqual, 1)

		//h.Dump()

		val, ok := h.LookupStringKey("hello")
		cv.So(val.GetString(), cv.ShouldStartWith, "world")
		cv.So(ok, cv.ShouldEqual, true)
	})

}
