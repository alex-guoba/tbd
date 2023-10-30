package provider

import (
	"sync"
)

type Agent struct {
	mu     sync.Mutex
	subs   map[string][]chan interface{}
	quit   chan struct{}
	closed bool
}

func NewAgent() *Agent {
	return &Agent{
		subs: make(map[string][]chan interface{}),
		quit: make(chan struct{}),
	}
}

func (b *Agent) Publish(topic string, msg interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	for _, ch := range b.subs[topic] {
		ch <- msg
	}
}

func (b *Agent) Subscribe(topic string) <-chan interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan interface{})
	b.subs[topic] = append(b.subs[topic], ch)
	return ch
}

func (b *Agent) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	close(b.quit)

	for _, ch := range b.subs {
		for _, sub := range ch {
			close(sub)
		}
	}
}
