package domain

import (
	"time"

	"github.com/google/uuid"
)

// ── Farm Entity ───────────────────────────────────────────────────────────────

type Farm struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Address     *string   `json:"address,omitempty"`
	Latitude    *float64  `json:"latitude,omitempty"`
	Longitude   *float64  `json:"longitude,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	SyncVersion int64     `json:"sync_version"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
