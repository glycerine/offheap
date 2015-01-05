package offheap_test

import (
	"encoding/binary"
	"fmt"

	"github.com/glycerine/offheap"
)

// compare for correctness checking
func hashEqualsMap(h *offheap.HashTable, m map[uint64]int) bool {
	if h.Population != uint64(len(m)) {
		fmt.Printf("h has size %d, but m has size %d.\n", h.Population, len(m))
		return false
	}
	var cell *offheap.Cell
	for k, v := range m {
		cell = h.Lookup(k)
		if cell == nil {
			fmt.Printf("m had key '%v', but h did not.\n", k)
			return false
		}
		if cell.GetInt() != v {
			fmt.Printf("m had key '%v':value '%v', but for that key, h had a different value: '%v'.\n", k, v, cell.Value)
			return false
		}
	}
	return true
}

func stringHashEqualsMap(h *offheap.StringHashTable, m map[string]int) bool {
	if h.Population != uint64(len(m)) {
		fmt.Printf("h has size %d, but m has size %d.\n", h.Population, len(m))
		return false
	}
	for k, v := range m {
		val, ok := h.LookupStringKey(k)
		if !ok {
			fmt.Printf("m had key '%v', but h did not.\n", k)
			return false
		}

		n := int(binary.LittleEndian.Uint64(val[:8]))

		if n != v {
			fmt.Printf("m had key '%v':value '%v', but for that key, h had a different value: '%v'.\n", k, v, n)
			return false
		}
	}
	return true
}
