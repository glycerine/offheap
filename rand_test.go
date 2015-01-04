package offheap_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/glycerine/go-offheap-hashtable"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestRandomOperationsOrder(t *testing.T) {

	h := offheap.NewHashTable(2)

	m := make(map[uint64]int)

	cv.Convey("given a sequence of random operations, the result should match what Go's builtin map does", t, func() {
		cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

		// basic insert
		m[1] = 2
		h.InsertIntValue(1, 2)
		cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

		m[3] = 4
		h.InsertIntValue(3, 4)
		cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

		// basic delete
		delete(m, 1)
		h.DeleteKey(1)
		cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

		delete(m, 3)
		h.DeleteKey(3)
		cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

		// now do random operations
		N := 1000
		seed := time.Now().UnixNano()
		fmt.Printf("\n TestRandomOperationsOrder() using seed = '%v'\n", seed)
		gen := rand.New(rand.NewSource(seed))

		for i := 0; i < N; i++ {

			op := gen.Int() % 4
			k := uint64(gen.Int() % (N / 4))
			v := gen.Int() % (N / 4)

			switch op {
			case 0, 1, 2:
				h.InsertIntValue(k, v)
				m[k] = v
				cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)
			case 3:
				h.DeleteKey(uint64(k))
				delete(m, k)
				cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

			}
		}

		// distribution more emphasizing deletes

		for i := 0; i < N; i++ {

			op := gen.Int() % 2
			k := uint64(gen.Int() % (N / 5))
			v := gen.Int() % (N / 2)

			switch op {
			case 0:
				h.InsertIntValue(k, v)
				m[k] = v
				cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)
			case 1:
				h.DeleteKey(uint64(k))
				delete(m, k)
				cv.So(offheap.HashEqualsMap(h, m), cv.ShouldEqual, true)

			}
		}
	})
}
