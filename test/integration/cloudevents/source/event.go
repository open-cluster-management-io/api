package source

import (
	"context"
	"sync"
)

// EventClient is a client that can watch resource status change events.
type EventClient struct {
	namespace string
	resChan   chan *Resource
}

func NewEventClient(namespace string) *EventClient {
	return &EventClient{
		namespace: namespace,
		resChan:   make(chan *Resource),
	}
}

func (c *EventClient) Receive() <-chan *Resource {
	return c.resChan
}

type EventHub struct {
	// lock for registered watch clients
	mu sync.RWMutex

	// Registered clients.
	clients map[*EventClient]struct{}

	// Inbound messages from the clients.
	broadcast chan *Resource
}

// NewEventHub creates a new EventHub.
func NewEventHub() *EventHub {
	return &EventHub{
		clients:   make(map[*EventClient]struct{}),
		broadcast: make(chan *Resource),
	}
}

func (h *EventHub) Register(client *EventClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = struct{}{}
}

func (h *EventHub) Unregister(client *EventClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, client)
	close(client.resChan)
}

func (h *EventHub) Broadcast(res *Resource) {
	h.broadcast <- res
}

func (h *EventHub) Update(res *Resource) {
	h.Broadcast(res)
}

func (h *EventHub) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				if client.namespace == res.Namespace || client.namespace == "+" {
					client.resChan <- res
				}
			}
			h.mu.RUnlock()
		}
	}
}
