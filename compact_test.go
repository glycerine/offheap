package offheap_test

import (
	"testing"

	"github.com/glycerine/go-offheap-hashtable"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestCompact(t *testing.T) {

	cv.Convey("Given a big table with no values it, Compact() should re-allocate it to be smaller", t, func() {
		h := offheap.NewHashTable(128)
		cv.So(len(h.Cells), cv.ShouldEqual, 128)
		cv.So(h.Population, cv.ShouldEqual, 0)

		h.Compact()

		cv.So(len(h.Cells), cv.ShouldEqual, 1)
		cv.So(h.Population, cv.ShouldEqual, 0)
	})

}
