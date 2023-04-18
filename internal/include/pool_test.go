package include_test

import (
	"fmt"
	"gosparse/internal/include"
	"testing"
)

func BenchmarkPool(b *testing.B) {
	b.Run("using pool", func(b *testing.B) {
		pool := include.UseIncPool(func() *include.Includes {
			return include.New()
		})

		for i := 0; i < b.N; i++ {
			inc := pool.Get()
			k := fmt.Sprintf("key-%d", i)
			inc.AddRel(k)
			pool.Release(inc)
		}
	})

	b.Run("without pool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			inc := include.New()
			k := fmt.Sprintf("key-%d", i)
			inc.AddRel(k)
		}
	})
}

// goos: linux
// goarch: amd64
// pkg: gosparse/internal/include
// cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
// === RUN   BenchmarkPool
//
// BenchmarkPool
//
// === RUN   BenchmarkPool/using_pool
// BenchmarkPool/using_pool
// BenchmarkPool/using_pool-12              1295466               858.0 ns/op           290 B/op          5 allocs/op
//
// === RUN   BenchmarkPool/without_pool
// BenchmarkPool/without_pool
// BenchmarkPool/without_pool-12            5122704               231.1 ns/op           224 B/op          5 allocs/op
//
// PASS
// ok      gosparse/internal/include       3.481s
