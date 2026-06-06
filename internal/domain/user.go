package domain

import (
	"time"

	"github.com/google/uuid"
)

// ── Role ──────────────────────────────────────────────────────────────────────

type Role string

const (
	RoleOwner  Role = "owner"
	RoleVet    Role = "vet"
	RoleWorker Role = "worker"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleVet, RoleWorker:
		return true
	}
	return false
}

// ── Tenant ────────────────────────────────────────────────────────────────────

type Tenant struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Plan              string    `json:"plan"`
	GlobalSyncVersion int64     `json:"global_sync_version"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}

// ── User ──────────────────────────────────────────────────────────────────────

type User struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never serialize
	FullName     string    `json:"full_name"`
	Role         Role      `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// ── Device ────────────────────────────────────────────────────────────────────

type Device struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	UserID          uuid.UUID  `json:"user_id"`
	DeviceUID       string     `json:"device_uid"`
	LastSyncVersion int64      `json:"last_sync_version"`
	LastSeenAt      *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
