package offheap_test

import (
	"testing"

	"github.com/glycerine/go-offheap-hashtable"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestIterator(t *testing.T) {

	cv.Convey("Given a table with 0,1,2 in it, the Iterator should give all three values back", t, func() {
		h := offheap.NewHashTable(8)
		cv.So(h.Population, cv.ShouldEqual, 0)
		for i := 0; i < 3; i++ {
			_, ok := h.Insert(uint64(i))
			cv.So(ok, cv.ShouldEqual, true)
			if i == 0 {
				// iterator should start with the zero value
				it := offheap.NewIterator(h)
				cv.So(it.Cur.Key, cv.ShouldEqual, 0)
			}
		}
		cv.So(h.Population, cv.ShouldEqual, 3)

		found := []uint64{}
		for it := offheap.NewIterator(h); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Key)
		}
		cv.So(len(found), cv.ShouldEqual, 3)
		cv.So(found, cv.ShouldContain, 0)
		cv.So(found, cv.ShouldContain, 1)
		cv.So(found, cv.ShouldContain, 2)
	})

	cv.Convey("Given a table with 1,2,3 in it, the Iterator should give all three values back", t, func() {
		h := offheap.NewHashTable(8)
		cv.So(h.Population, cv.ShouldEqual, 0)
		for i := 1; i < 4; i++ {
			_, ok := h.Insert(uint64(i))
			cv.So(ok, cv.ShouldEqual, true)
			if i == 0 {
				// iterator should not start with the zero value, not inserted.
				it := offheap.NewIterator(h)
				cv.So(it.Cur.Key, cv.ShouldEqual, 1)
			}
		}
		cv.So(h.Population, cv.ShouldEqual, 3)

		found := []uint64{}
		for it := offheap.NewIterator(h); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Key)
		}
		cv.So(len(found), cv.ShouldEqual, 3)
		cv.So(found, cv.ShouldContain, 3)
		cv.So(found, cv.ShouldContain, 1)
		cv.So(found, cv.ShouldContain, 2)
	})

}
