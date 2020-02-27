package gopubsub

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// PubSub implements publish/subscribe messaging paradigm
type PubSub interface {
	// Pub publishes arguments to the given topic subscribers
	// Pub block only when the buffer of one of the subscribers is full.
	Pub(topic string, args ...interface{}) error
	// Close unsubscribe all handlers from given topic
	Close(topic string)
	// Sub subscribes to the given topic
	Sub(topic string, fn interface{}) error
	// Unsub unsubscribe handlerImpl from the given topic
	Unsub(topic string, fn interface{}) error
}

type handlersMap map[string][]handler

type handler interface {
	pub([]reflect.Value)
	close()
	isCallback(f reflect.Value) bool
}

type handlerImpl struct {
	callback reflect.Value
	queue    chan []reflect.Value
}

func makeHandle(fn interface{}, handlerQueueSize int) *handlerImpl {
	h := &handlerImpl{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, handlerQueueSize),
	}

	go h.start()

	return h
}

func (h *handlerImpl) pub(v []reflect.Value)           { h.queue <- v }
func (h *handlerImpl) close()                          { close(h.queue) }
func (h *handlerImpl) isCallback(f reflect.Value) bool { return h.callback == f }

func (h *handlerImpl) start() {
	for args := range h.queue {
		h.callback.Call(args)
	}
}

type messageBus struct {
	handlerQueueSize int
	mtx              sync.RWMutex
	handlers         handlersMap
}

// ErrTopicNotExists is the error to indicate that the topic does not exist.
var ErrTopicNotExists = errors.New("topic does not exist")

func (b *messageBus) Pub(topic string, args ...interface{}) error {
	rArgs := buildHandlerArgs(args)

	b.mtx.RLock()
	defer b.mtx.RUnlock()

	hs, ok := b.handlers[topic]

	if !ok {
		return fmt.Errorf("topic %s doesn't exist %w", topic, ErrTopicNotExists)
	}

	for _, h := range hs {
		h.pub(rArgs)
	}

	return nil
}

func (b *messageBus) Sub(topic string, fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	h := makeHandle(fn, b.handlerQueueSize)

	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.handlers[topic] = append(b.handlers[topic], h)

	return nil
}

func (b *messageBus) Unsub(topic string, fn interface{}) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	handlers, ok := b.handlers[topic]
	if !ok {
		return fmt.Errorf("topic %s doesn't exist %w", topic, ErrTopicNotExists)
	}

	rv := reflect.ValueOf(fn)

	for i, h := range handlers {
		if h.isCallback(rv) {
			h.close()

			b.handlers[topic] = append(handlers[:i], handlers[i+1:]...)
		}
	}

	if len(b.handlers[topic]) == 0 {
		delete(b.handlers, topic)
	}

	return nil
}

func (b *messageBus) Close(topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	handlers, ok := b.handlers[topic]
	if !ok {
		return
	}

	for _, h := range handlers {
		h.close()
	}

	delete(b.handlers, topic)
}

func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, len(args))

	for i, arg := range args {
		reflectedArgs[i] = reflect.ValueOf(arg)
	}

	return reflectedArgs
}

// New creates new PubSub
// handlerQueueSize sets buffered channel length per subscriber
func New(handlerQueueSize int) PubSub {
	if handlerQueueSize == 0 {
		panic("handlerQueueSize has to be greater then 0")
	}

	return &messageBus{
		handlerQueueSize: handlerQueueSize,
		handlers:         make(handlersMap),
	}
}
