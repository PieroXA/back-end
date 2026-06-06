package domain

import (
	"time"

	"github.com/google/uuid"
)

// ── Movement Reason ───────────────────────────────────────────────────────────

type MovementReason string

const (
	ReasonSaleTransfer MovementReason = "SALE_TRANSFER"
	ReasonGrazing      MovementReason = "GRAZING"
	ReasonVeterinary   MovementReason = "VETERINARY"
	ReasonReturn       MovementReason = "RETURN"
	ReasonPurchase     MovementReason = "PURCHASE"
	ReasonOther        MovementReason = "OTHER"
)

func (r MovementReason) IsValid() bool {
	switch r {
	case ReasonSaleTransfer, ReasonGrazing, ReasonVeterinary, ReasonReturn, ReasonPurchase, ReasonOther:
		return true
	}
	return false
}

// ── Movement Entity ───────────────────────────────────────────────────────────

type Movement struct {
	ID           uuid.UUID      `json:"id"`
	TenantID     uuid.UUID      `json:"tenant_id"`
	AnimalID     uuid.UUID      `json:"animal_id"`
	FromFarmID   *uuid.UUID     `json:"from_farm_id,omitempty"`
	ToFarmID     *uuid.UUID     `json:"to_farm_id,omitempty"`
	MovementDate time.Time      `json:"movement_date"`
	Reason       MovementReason `json:"reason"`
	Notes        *string        `json:"notes,omitempty"`
	SyncVersion  int64          `json:"sync_version"`
	IsDeleted    bool           `json:"is_deleted"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
