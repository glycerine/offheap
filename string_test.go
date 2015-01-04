package offheap_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/glycerine/go-offheap-hashtable"
	cv "github.com/smartystreets/goconvey/convey"
)

type Num struct {
	name string
	num  int
}

func TestRandomStringOps(t *testing.T) {

	// make a reference array that maps i -> hex string version,
	// for testing the string <-> int hashmap implementation
	M := 36
	nm := make([]Num, M)
	for i := 0; i < M; i++ {
		nm[i].name = strconv.FormatInt(int64(i), 36)
		nm[i].num = i
	}

	h := offheap.NewHashTable(2)

	m := make(map[string]int)

	cv.Convey("given a sequence of random operations, the result should match what Go's builtin map does", t, func() {
		cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

		// basic insert
		m[nm[0].name] = nm[0].num
		h.InsertStringKey(nm[0].name, nm[0].num)
		cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

		m[nm[1].name] = nm[1].num
		h.InsertStringKey(nm[1].name, nm[1].num)
		cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

		// basic delete
		delete(m, nm[0].name)
		h.DeleteStringKey(nm[0].name)
		cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

		delete(m, nm[1].name)
		h.DeleteStringKey(nm[1].name)
		cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

		// now do random operations
		N := 1000
		seed := time.Now().UnixNano()
		fmt.Printf("\n TestRandomOperationsOrder() using seed = '%v'\n", seed)
		gen := rand.New(rand.NewSource(seed))

		for i := 0; i < N; i++ {

			op := gen.Int() % 4
			w := uint64(gen.Int() % M)

			switch op {
			case 0, 1, 2:
				m[nm[w].name] = nm[w].num
				h.InsertStringKey(nm[w].name, nm[w].num)
				cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

			case 3:
				delete(m, nm[w].name)
				h.DeleteStringKey(nm[w].name)
				cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)
			}
		}

		// distribution more emphasizing deletes

		for i := 0; i < N; i++ {

			op := gen.Int() % 2
			w := uint64(gen.Int() % M)

			switch op {
			case 0:
				m[nm[w].name] = nm[w].num
				h.InsertStringKey(nm[w].name, nm[w].num)
				cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

			case 1:
				delete(m, nm[w].name)
				h.DeleteStringKey(nm[w].name)
				cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)
			}
		}

		fmt.Printf("\n === h at the end of %d string ops:\n", N)
		h.DumpStringKey()
	})
}
