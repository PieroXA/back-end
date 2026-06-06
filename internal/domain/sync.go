package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ── Sync Operation (received from client Push) ────────────────────────────────

type SyncOperation struct {
	OpID            string          `json:"op_id"`
	EntityType      string          `json:"entity_type"` // animal|health_event|movement|farm
	EntityID        string          `json:"entity_id"`
	Action          string          `json:"action"` // CREATE|UPDATE|DELETE
	Payload         json.RawMessage `json:"payload"`
	ChangedFields   []string        `json:"changed_fields"`
	ClientTimestamp time.Time       `json:"client_timestamp"`
}

// ── Sync Push Request / Response ──────────────────────────────────────────────

type SyncPushRequest struct {
	DeviceID        string          `json:"device_id"`
	LastSyncVersion int64           `json:"last_sync_version"`
	Operations      []SyncOperation `json:"operations"`
}

type SyncPushResponse struct {
	Applied        int              `json:"applied"`
	Conflicts      []ConflictResult `json:"conflicts"`
	NewSyncVersion int64            `json:"new_sync_version"`
}

type ConflictResult struct {
	OpID       string `json:"op_id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Field      string `json:"field"`
	ClientVal  any    `json:"client_value"`
	ServerVal  any    `json:"server_value"`
	ClientTime string `json:"client_timestamp"`
	ServerTime string `json:"server_timestamp"`
}

// ── Pull Response ─────────────────────────────────────────────────────────────

type PullResponse struct {
	Animals      []*Animal      `json:"animals"`
	HealthEvents []*HealthEvent `json:"health_events"`
	Movements    []*Movement    `json:"movements"`
	Farms        []*Farm        `json:"farms"`
	SyncVersion  int64          `json:"sync_version"`
}

// ── Notification ──────────────────────────────────────────────────────────────

type Notification struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	HealthEventID *uuid.UUID `json:"health_event_id,omitempty"`
	Type          string     `json:"type"`
	Message       string     `json:"message"`
	IsRead        bool       `json:"is_read"`
	ScheduledDate time.Time  `json:"scheduled_date"`
	CreatedAt     time.Time  `json:"created_at"`
}
