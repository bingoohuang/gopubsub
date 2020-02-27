// nolint gomnd
package gopubsub_test

import (
	"fmt"
	"sync"

	"github.com/bingoohuang/gopubsub"
)

func Example() {
	queueSize := 100
	bus := gopubsub.New(queueSize)

	var wg sync.WaitGroup

	wg.Add(2)

	_ = bus.Sub("topic", func(v bool) {
		defer wg.Done()
		fmt.Println("s1", v)
	})

	_ = bus.Sub("topic", func(v bool) {
		defer wg.Done()
		fmt.Println("s2", v)
	})

	// Pub block only when the buffer of one of the subscribers is full.
	// change the buffer size altering queueSize when creating new gopubsub
	bus.Pub("topic", true)
	wg.Wait()

	// Unordered output:
	// s1 true
	// s2 true
}

func Example_second() {
	queueSize := 2
	subscribersAmount := 3

	ch := make(chan int, queueSize)
	defer close(ch)

	bus := gopubsub.New(queueSize)

	for i := 0; i < subscribersAmount; i++ {
		_ = bus.Sub("topic", func(i int, out chan<- int) { out <- i })
	}

	go func() {
		for n := 0; n < queueSize; n++ {
			bus.Pub("topic", n, ch)
		}
	}()

	var sum = 0
	for sum < (subscribersAmount * queueSize) {
		select {
		case <-ch:
			sum++
		}
	}

	fmt.Println(sum)
	// Output:
	// 6
}
