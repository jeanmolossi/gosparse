package include_test

import (
	"fmt"
	"testing"
)

const LIST_SIZE = 5

var (
	mapped  = make(map[string]struct{}, LIST_SIZE)
	arrayed = make([]string, LIST_SIZE)
)

func initialize() {
	for i := 0; i < LIST_SIZE; i++ {
		key := fmt.Sprintf("key-%d", i)
		mapped[key] = struct{}{}
		arrayed = append(arrayed, key)
	}
}

func searchMap(k string) string {
	if _, exists := mapped[k]; exists {
		return k
	}

	return ""
}

func searchArray(k string) string {
	for _, v := range arrayed {
		if v == k {
			return v
		}
	}

	return ""
}

func BenchmarkSearch(b *testing.B) {
	initialize()

	b.Run("search map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			k := fmt.Sprintf("key-%d", i)
			if i > LIST_SIZE {
				k = fmt.Sprintf("key-%d", LIST_SIZE)
			}

			searchMap(k)
		}
	})

	b.Run("search array", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			k := fmt.Sprintf("key-%d", i)
			if i > LIST_SIZE {
				k = fmt.Sprintf("key-%d", LIST_SIZE)
			}

			searchArray(k)
		}
	})
}
