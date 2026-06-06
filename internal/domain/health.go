package domain

import (
	"time"

	"github.com/google/uuid"
)

// ── Event Type ────────────────────────────────────────────────────────────────

type EventType string

const (
	EventVaccine   EventType = "VACCINE"
	EventTreatment EventType = "TREATMENT"
	EventDiagnosis EventType = "DIAGNOSIS"
	EventDeworming EventType = "DEWORMING"
	EventSurgery   EventType = "SURGERY"
	EventCheckup   EventType = "CHECKUP"
)

func (e EventType) IsValid() bool {
	switch e {
	case EventVaccine, EventTreatment, EventDiagnosis, EventDeworming, EventSurgery, EventCheckup:
		return true
	}
	return false
}

// RequiresProduct returns true for event types that mandate a product_name.
func (e EventType) RequiresProduct() bool {
	return e == EventVaccine || e == EventTreatment || e == EventDeworming
}

// ── Health Event Entity ───────────────────────────────────────────────────────

type HealthEvent struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	AnimalID    uuid.UUID  `json:"animal_id"`
	EventType   EventType  `json:"event_type"`
	EventDate   time.Time  `json:"event_date"`
	ProductName *string    `json:"product_name,omitempty"`
	DoseMl      *float64   `json:"dose_ml,omitempty"`
	LotNumber   *string    `json:"lot_number,omitempty"`
	VetName     *string    `json:"vet_name,omitempty"`
	NextDate    *time.Time `json:"next_date,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
	Cost        *float64   `json:"cost,omitempty"`
	BulkBatchID *uuid.UUID `json:"bulk_batch_id,omitempty"`
	SyncVersion int64      `json:"sync_version"`
	IsDeleted   bool       `json:"is_deleted"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
