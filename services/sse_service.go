package services

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

// Event represents a generic SSE event
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// SSEService manages SSE connections and broadcasting
type SSEService struct {
	// Map of communityID -> list of client channels
	clients map[uuid.UUID][]chan Event
	// Mutex to protect the clients map
	mu sync.RWMutex
}

// NewSSEService creates a new SSEService
func NewSSEService() *SSEService {
	return &SSEService{
		clients: make(map[uuid.UUID][]chan Event),
	}
}

// Subscribe adds a new client to a community's stream
// Returns a channel that receives events
func (s *SSEService) Subscribe(communityID uuid.UUID) chan Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Buffer of 10 to prevent blocking if client is slow
	clientChan := make(chan Event, 10)
	s.clients[communityID] = append(s.clients[communityID], clientChan)

	log.Printf("Client subscribed to community %s. Total clients: %d", communityID, len(s.clients[communityID]))
	return clientChan
}

// Unsubscribe removes a client from a community's stream
func (s *SSEService) Unsubscribe(communityID uuid.UUID, clientChan chan Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := s.clients[communityID]
	for i, c := range clients {
		if c == clientChan {
			// Remove from slice
			s.clients[communityID] = append(clients[:i], clients[i+1:]...)
			close(clientChan)
			log.Printf("Client unsubscribed from community %s. Total clients: %d", communityID, len(s.clients[communityID]))
			break
		}
	}

	// Clean up map entry if empty
	if len(s.clients[communityID]) == 0 {
		delete(s.clients, communityID)
	}
}

// Broadcast sends an event to all clients in a community
func (s *SSEService) Broadcast(communityID uuid.UUID, eventType string, payload interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients, ok := s.clients[communityID]
	if !ok {
		return
	}

	event := Event{
		Type:    eventType,
		Payload: payload,
	}

	for _, clientChan := range clients {
		// Non-blocking send to avoid holding up the broadcaster
		select {
		case clientChan <- event:
		default:
			log.Printf("Slow client in community %s, dropping event", communityID)
		}
	}

	if len(clients) > 0 {
		log.Printf("Broadcasted event '%s' to %d clients in community %s", eventType, len(clients), communityID)
	}
}

// BroadcastJSON is a helper to broadcast raw JSON if needed, but we use typed Event struct above
func (s *SSEService) BroadcastPostCreated(communityID, postID uuid.UUID, post interface{}) {
	s.Broadcast(communityID, "post_created", post)
}
