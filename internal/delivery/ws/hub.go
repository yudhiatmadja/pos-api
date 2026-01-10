package ws

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Redis Integration
	RedisClient *redis.Client
	PubSub      *redis.PubSub
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[*Client]bool),
		RedisClient: rdb,
	}
}

func (h *Hub) Run() {
	// Subscribe to global order events channel
	// In production, might subscribe to specific outlet channels
	h.PubSub = h.RedisClient.Subscribe(context.Background(), "order_updates")
	ch := h.PubSub.Channel()

	go func() {
		for msg := range ch {
			h.broadcast([]byte(msg.Payload))
		}
	}()

	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		}
	}
}

func (h *Hub) broadcast(message []byte) {
	// Naive broadcast to all connected clients (KDS, POS)
	// Filter logic should be added here based on OutletID
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

// PublishEvent allows other parts of app to publish events
func (h *Hub) PublishEvent(ctx context.Context, eventType string, payload interface{}) error {
	msg := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// Publish to Redis
	return h.RedisClient.Publish(ctx, "order_updates", bytes).Err()
}
