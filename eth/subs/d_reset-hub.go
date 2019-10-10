package subs

import (
	"sync"

	"github.com/phonkee/go-pubsub"
)

// LegacyEthResetHub implements Hub interface
// Remembers which subscribers messages have been published to
// and does not send repeat messages to any subscribers.
// Revert resets the memory of the subscribers that have received messages.
type LegacyEthResetHub struct {
	mutex    *sync.RWMutex
	registry map[pubsub.Subscriber]bool
}

// New returns new hub instance. hub is goroutine safe.
func NewLegacyEthResetHub() pubsub.ResetHub {
	return &LegacyEthResetHub{
		mutex:    &sync.RWMutex{},
		registry: map[pubsub.Subscriber]bool{},
	}
}

// CloseSubscriber removes subscriber from hub
func (h *LegacyEthResetHub) CloseSubscriber(subscriber pubsub.Subscriber) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.registry, subscriber)
}

// Publish publishes message to subscribers
// todo Warning this function can throw an exception
func (h *LegacyEthResetHub) Publish(message pubsub.Message) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	count := 0
	// iterate over all subscribers, and publish messages
	for sub := range h.registry {
		if sub.Match(message.Topic()) {
			count += sub.Publish(message)
		}
	}

	return count
}

// Subscribe adds subscription to topics and returns subscriber
func (h *LegacyEthResetHub) Subscribe(topics ...string) pubsub.Subscriber {
	var result pubsub.Subscriber
	if len(topics) > 0 {
		result = newEthSubscriber(h, topics[0])
	} else {
		result = newEthSubscriber(h, "")
	}

	h.mutex.Lock()
	h.registry[result] = true
	h.mutex.Unlock()

	return result
}

func (h *LegacyEthResetHub) Reset() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for sub := range h.registry {
		h.registry[sub] = true
	}
}
