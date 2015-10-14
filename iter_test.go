package offheap

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIterator(t *testing.T) {

	Convey("Given a table with 0,1,2 in it, the Iterator should give all three values back", t, func() {
		h := NewHashTableInt(8)
		for i := 0; i < 3; i++ {
			h.InsertInt(uint64(i), i+10)
			if i == 0 {
				// iterator should start with the zero value
				it := h.NewIterator()
				So(it.Cur.Value, ShouldEqual, 10)
			}
		}
		So(h.Population, ShouldEqual, 3)

		found := []int{}
		for it := h.NewIterator(); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Value)
		}
		So(len(found), ShouldEqual, 3)
		So(found, ShouldContain, 10)
		So(found, ShouldContain, 11)
		So(found, ShouldContain, 12)
	})

	Convey("Given a table with 1,2,3 in it, the Iterator should give all three values back", t, func() {
		h := NewHashTableInt(8)
		So(h.Population, ShouldEqual, 0)
		for i := 1; i < 4; i++ {
			h.InsertInt(uint64(i), i+100)
			if i == 1 {
				// iterator should not start with the zero value, not inserted.
				it := h.NewIterator()
				So(it.Cur, ShouldNotBeNil)
				So(it.Cur.Value, ShouldEqual, 101)
			}
		}
		So(h.Population, ShouldEqual, 3)

		found := []int{}
		for it := h.NewIterator(); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Value)
		}
		So(len(found), ShouldEqual, 3)
		So(found, ShouldContain, 101)
		So(found, ShouldContain, 102)
		So(found, ShouldContain, 103)
	})

	Convey("Given a table with the regular 0-th slot and the special zero-location spot occupied, then the the Iterator should still give all the values back", t, func() {
		h := NewHashTableInt(4)
		for i := 0; i < 2; i++ {
			h.InsertInt(uint64(i), i+200)
			if i == 0 {
				// iterator should start with the zero value
				it := h.NewIterator()
				So(it.Cur.Value, ShouldEqual, 200)
			}
		}
		So(h.Population, ShouldEqual, 2)

		found := []int{}
		for it := h.NewIterator(); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Value)
		}
		So(len(found), ShouldEqual, 2)
		So(found, ShouldContain, 200)
		So(found, ShouldContain, 201)
	})

	Convey("Given a table with the regular 0-th slot *empty* and the special zero-location spot occupied, then the the Iterator should still give all the values back", t, func() {
		h := NewHashTableInt(8)
		So(h.Population, ShouldEqual, 0)
		for i := 0; i < 2; i++ {
			h.InsertInt(uint64(i), 300+i)
			if i == 0 {
				// iterator should start with the zero value
				it := h.NewIterator()
				So(it.Cur.Value, ShouldEqual, 300)
				So(it.Pos, ShouldEqual, -1)
			}
		}
		So(h.Population, ShouldEqual, 2)

		found := []int{}
		for it := h.NewIterator(); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Value)
		}
		So(len(found), ShouldEqual, 2)
		So(found, ShouldContain, 300)
		So(found, ShouldContain, 301)
	})

	Convey("Given a table with the regular 0-th slot *filled* and the special zero-location spot *empty*, then the the Iterator should still give the one value back", t, func() {
		// size 4 just happens to generate an occupation at h.Cells[0]
		h := NewHashTableInt(4)
		i := 1
		h.InsertInt(uint64(i), 1001)

		// iterator should start with the zero value
		it := h.NewIterator()
		So(it.Cur.Value, ShouldEqual, 1001)
		So(it.Pos, ShouldEqual, 0)

		So(h.Population, ShouldEqual, 1)

		found := []int{}
		for it := h.NewIterator(); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Value)
		}
		So(len(found), ShouldEqual, 1)
		So(found, ShouldContain, 1001)
	})

	Convey("Given an empty table, an Iterator should still work fine, without crashing", t, func() {
		h := NewHashTableInt(4)
		So(h.Population, ShouldEqual, 0)

		found := []int{}
		for it := h.NewIterator(); it.Cur != nil; it.Next() {
			found = append(found, it.Cur.Value)
		}
		So(len(found), ShouldEqual, 0)
	})

}
