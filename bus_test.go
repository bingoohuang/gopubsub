// nolint gomnd
package gopubsub

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	bus := New(runtime.NumCPU())

	if bus == nil {
		t.Fail()
	}
}
func TestNewError(t *testing.T) {
	defer func() {
		if err := recover().(string); err != "handlerQueueSize has to be greater then 0" {
			t.Fatalf("Wrong panic message: %s", err)
		}
	}()

	New(0)
}

func TestCloseNoSubTopic(t *testing.T) {
	bus := New(1)
	bus.Pub("na", func() {})
	bus.Close("na")
}

func TestSubscribe(t *testing.T) {
	bus := New(runtime.NumCPU())

	if bus.Sub("test", func() {}) != nil {
		t.Fail()
	}

	if bus.Sub("test", 2) == nil {
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New(runtime.NumCPU())

	handler := func() {}

	if err := bus.Sub("test", handler); err != nil {
		t.Fatal(err)
	}

	if err := bus.Unsub("test", handler); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if err := bus.Unsub("non-existed", func() {}); err == nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestClose(t *testing.T) {
	bus := New(runtime.NumCPU())

	handler := func() {}

	if err := bus.Sub("test", handler); err != nil {
		t.Fatal(err)
	}

	original, ok := bus.(*pubSub)
	if !ok {
		fmt.Println("Could not cast message bus to its original type")
		t.Fail()
	}

	if len(original.handlers) == 0 {
		fmt.Println("Did not subscribed handlerImpl to topic")
		t.Fail()
	}

	bus.Close("test")

	if len(original.handlers) != 0 {
		fmt.Println("Did not unsubscribed handlers from topic")
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	bus := New(runtime.NumCPU())

	var wg sync.WaitGroup

	wg.Add(2)

	first := false
	second := false

	if err := bus.Sub("topic", func(v bool) {
		defer wg.Done()
		first = v
	}); err != nil {
		t.Fatal(err)
	}

	if err := bus.Sub("topic", func(v bool) {
		defer wg.Done()
		second = v
	}); err != nil {
		t.Fatal(err)
	}

	_ = bus.Pub("topic", true)

	wg.Wait()

	if first == false || second == false {
		t.Fail()
	}
}

func TestHandleError(t *testing.T) {
	bus := New(runtime.NumCPU())
	if err := bus.Sub("topic", func(out chan<- error) {
		out <- errors.New("I do throw error")
	}); err != nil {
		t.Fatal(err)
	}

	out := make(chan error)
	defer close(out)

	_ = bus.Pub("topic", out)

	if <-out == nil {
		t.Fail()
	}
}
