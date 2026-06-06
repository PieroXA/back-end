package animal

import (
	"time"

	"github.com/google/uuid"
)

// CreateAnimalDTO holds data required to create an animal.
type CreateAnimalDTO struct {
	TenantID  uuid.UUID
	FarmID    uuid.UUID
	VisualTag string
	RfidTag   *string
	Species   string
	Breed     *string
	Sex       string
	BirthDate time.Time
	WeightKg  *float64
	Color     *string
	SireID    *uuid.UUID
	DamID     *uuid.UUID
	Notes     *string
}

// UpdateAnimalDTO holds data for semantic updates (PATCH).
type UpdateAnimalDTO struct {
	TenantID  uuid.UUID
	AnimalID  uuid.UUID
	VisualTag *string
	RfidTag   *string
	Breed     *string
	WeightKg  *float64
	Color     *string
	SireID    *uuid.UUID
	DamID     *uuid.UUID
	FarmID    *uuid.UUID // Can trigger a Movement creation
	Notes     *string
}

// ChangeStatusDTO holds data to change an animal's lifecycle status.
type ChangeStatusDTO struct {
	TenantID    uuid.UUID
	AnimalID    uuid.UUID
	Status      string
	StatusDate  time.Time
	StatusNotes *string
	BuyerName   *string // For 'sold' status
	CauseDeath  *string // For 'dead' status
}
