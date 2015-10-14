package offheap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// compare for correctness checking
func hashEqualsMap(h *HashTableInt, m map[uint64]int) bool {
	if h.Population != uint64(len(m)) {
		fmt.Printf("h has size %d, but m has size %d.\n", h.Population, len(m))
		return false
	}
	for k, v := range m {
		cell := h.Lookup(k)
		if cell == nil {
			fmt.Printf("m had key '%v', but h did not.\n", k)
			return false
		}
		if cell.Value != v {
			fmt.Printf("m had key '%v':value '%v', but for that key, h had a different value: '%v'.\n", k, v, cell.Value)
			return false
		}
	}
	return true
}

func TestRandomOperationsOrder(t *testing.T) {

	h := NewHashTableInt(2)

	m := make(map[uint64]int)

	Convey("given a sequence of random operations, the result should match what Go's builtin map does", t, func() {
		So(hashEqualsMap(h, m), ShouldEqual, true)

		// basic insert
		m[1] = 2
		h.InsertInt(1, 2)

		// h.InsertIntValue(1, 2)
		So(hashEqualsMap(h, m), ShouldEqual, true)

		m[3] = 4
		h.InsertInt(3, 4)
		So(hashEqualsMap(h, m), ShouldEqual, true)

		// basic delete
		delete(m, 1)
		h.DeleteKey(1)
		So(hashEqualsMap(h, m), ShouldEqual, true)

		delete(m, 3)
		h.DeleteKey(3)
		So(hashEqualsMap(h, m), ShouldEqual, true)

		// now do random operations
		N := 1000
		seed := time.Now().UnixNano()
		gen := rand.New(rand.NewSource(seed))

		for i := 0; i < N; i++ {

			op := gen.Int() % 4
			k := uint64(gen.Int() % (N / 4))
			v := gen.Int() % (N / 4)

			switch op {
			case 0, 1, 2:
				h.InsertInt(k, v)
				m[k] = v
				So(hashEqualsMap(h, m), ShouldEqual, true)
			case 3:
				h.DeleteKey(uint64(k))
				delete(m, k)
				So(hashEqualsMap(h, m), ShouldEqual, true)

			}
		}

		// distribution more emphasizing deletes

		for i := 0; i < N; i++ {

			op := gen.Int() % 2
			k := uint64(gen.Int() % (N / 5))
			v := gen.Int() % (N / 2)

			switch op {
			case 0:
				h.InsertInt(k, v)
				m[k] = v
				So(hashEqualsMap(h, m), ShouldEqual, true)
			case 1:
				h.DeleteKey(uint64(k))
				delete(m, k)
				So(hashEqualsMap(h, m), ShouldEqual, true)

			}
		}
	})
}
