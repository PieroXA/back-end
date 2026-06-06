package websocket

import (
	"log"

	"mifinca-backend/internal/domain/ports"
)

// MockHub is a temporary mock for the EventBus.
// In Sprint 3, this will be replaced by the actual WebSocket Hub.
type MockHub struct{}

func NewMockHub() *MockHub {
	return &MockHub{}
}

func (h *MockHub) Publish(e ports.Event) {
	log.Printf("[EventBus] Broadcast: EventType=%s EntityType=%s EntityID=%s TenantID=%s\n",
		e.EventType, e.EntityType, e.EntityID, e.TenantID)
}
