package ports

import "encoding/json"

// Event represents a domain event to be broadcast via WebSocket to tenant devices.
type Event struct {
	EventType  string          `json:"event_type"`  // ANIMAL_CREATED, ANIMAL_UPDATED, HEALTH_CREATED, etc.
	EntityType string          `json:"entity_type"` // animal, health_event, movement, farm
	EntityID   string          `json:"entity_id"`
	TenantID   string          `json:"tenant_id"`
	Payload    json.RawMessage `json:"payload"`
	Timestamp  string          `json:"timestamp"` // RFC3339
}

// EventBus publishes events to all connected WebSocket clients of the tenant.
type EventBus interface {
	Publish(e Event)
}
