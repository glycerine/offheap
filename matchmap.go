package offheap

import "fmt"

// compare for correctness checking
func HashEqualsMap(h *HashTable, m map[uint64]int) bool {
	if h.Population != uint64(len(m)) {
		fmt.Printf("h has size %d, but m has size %d.\n", h.Population, len(m))
		return false
	}
	var cell *Cell
	for k, v := range m {
		cell = h.Lookup(k)
		if cell == nil {
			fmt.Printf("m had key '%v', but h did not.\n", k)
			return false
		}
		if cell.Value.(int) != v {
			fmt.Printf("m had key '%v':value '%v', but for that key, h had a different value: '%v'.\n", k, v, cell.Value.(int))
			return false
		}
	}
	return true
}

func StringHashEqualsMap(h *HashTable, m map[string]int) bool {
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
		if val.(int) != v {
			fmt.Printf("m had key '%v':value '%v', but for that key, h had a different value: '%v'.\n", k, v, val.(int))
			return false
		}
	}
	return true
}
