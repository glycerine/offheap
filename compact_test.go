package offheap

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCompact(t *testing.T) {

	Convey("Given a big table with no values it, Compact() should re-allocate it to be smaller", t, func() {
		h := NewHashTableInt(128)
		So(h.ArraySize, ShouldEqual, 128)
		So(h.Population, ShouldEqual, 0)

		h.Compact()

		So(h.ArraySize, ShouldEqual, 1)
		So(h.Population, ShouldEqual, 0)
	})

}

func TestCompatAfterDelete(t *testing.T) {

	h := NewHashTableInt(2)

	h.Clear()
	Convey("after being filled up and then deleted down to just 2 elements, Compact() should reduce table size to 4", t, func() {
		N := 100

		fmt.Println("Array is", h.ArraySize)

		// fill up a table
		So(h.Population, ShouldEqual, 0)
		for i := 0; i < N; i++ {
			h.Insert(uint64(i))
			So(h.Population, ShouldEqual, i+1)
		}

		// now delete from it
		for i := 0; i < N-2; i++ {
			h.DeleteKey(uint64(i))
		}
		So(h.ArraySize, ShouldEqual, 256)
		h.Compact()
		So(h.ArraySize, ShouldEqual, 4)

	})
}
