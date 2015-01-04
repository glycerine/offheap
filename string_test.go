package offheap_test

import (
	"strconv"
	"testing"

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
	M := 10
	nm := make([]Num, M)
	for i := 0; i < M; i++ {
		nm[i].name = strconv.FormatInt(int64(i), 16)
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

		/*
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
					cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)
				case 3:
					h.DeleteKey(uint64(k))
					delete(m, k)
					cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

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
					cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)
				case 1:
					h.DeleteKey(uint64(k))
					delete(m, k)
					cv.So(offheap.StringHashEqualsMap(h, m), cv.ShouldEqual, true)

				}
			}
		*/
	})
}
