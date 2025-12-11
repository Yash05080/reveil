package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSSEService(t *testing.T) {
	service := NewSSEService()
	communityID := uuid.New()

	// Test Subscription
	clientChan := service.Subscribe(communityID)
	if clientChan == nil {
		t.Fatal("Subscribe returned nil channel")
	}

	// Test Broadcast
	testPayload := map[string]string{"message": "hello"}
	go func() {
		service.Broadcast(communityID, "test_event", testPayload)
	}()

	// Wait for event
	select {
	case event := <-clientChan:
		if event.Type != "test_event" {
			t.Errorf("Expected event type 'test_event', got %s", event.Type)
		}
		// Verify payload
		payloadBytes, _ := json.Marshal(event.Payload)
		expectedBytes, _ := json.Marshal(testPayload)
		if string(payloadBytes) != string(expectedBytes) {
			t.Errorf("Expected payload %s, got %s", string(expectedBytes), string(payloadBytes))
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event")
	}

	// Test Unsubscribe
	service.Unsubscribe(communityID, clientChan)

	// Verify channel is closed
	select {
	case _, ok := <-clientChan:
		if ok {
			t.Error("Channel should be closed after unsubscribe")
		}
	case <-time.After(100 * time.Millisecond):
		// This is good, reading closed channel should return immediately
	}
}
