// nolint gomnd
package gopubsub

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkPublish(b *testing.B) {
	bus := New(b.N)

	var wg sync.WaitGroup

	wg.Add(b.N)

	_ = bus.Sub("topic", func() {
		wg.Done()
	})

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bus.Pub("topic")
		}
	})

	wg.Wait()
}

func BenchmarkSubscribe(b *testing.B) {
	bus := New(1)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bus.Sub("topic", func() {})
		}
	})
}

func benchmark(b *testing.B, subscribersCount, topicsCount int) {
	bus := New(b.N)

	var wg sync.WaitGroup

	wg.Add(b.N * subscribersCount)

	for i := 0; i < topicsCount; i++ {
		for j := 0; j < subscribersCount; j++ {
			_ = bus.Sub(strconv.Itoa(i), func() {
				wg.Done()
			})
		}
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bus.Pub(strconv.Itoa(topicsCount - 1))
		}
	})

	wg.Wait()
}

func Benchmark1Subscriber1Topic(b *testing.B) {
	benchmark(b, 1, 1)
}

func Benchmark1Subscriber100Topics(b *testing.B) {
	benchmark(b, 1, 100)
}

func Benchmark100Subscribers1Topic(b *testing.B) {
	benchmark(b, 100, 1)
}

func Benchmark100Subscribers100Topics(b *testing.B) {
	benchmark(b, 100, 100)
}
