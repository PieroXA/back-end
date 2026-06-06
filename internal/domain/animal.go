package domain

import (
	"time"

	"github.com/google/uuid"
)

// ── Animal Status ─────────────────────────────────────────────────────────────

type AnimalStatus string

const (
	StatusActive AnimalStatus = "active"
	StatusSold   AnimalStatus = "sold"
	StatusDead   AnimalStatus = "dead"
	StatusLost   AnimalStatus = "lost"
)

func (s AnimalStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusSold, StatusDead, StatusLost:
		return true
	}
	return false
}

func (s AnimalStatus) IsFinal() bool {
	return s == StatusSold || s == StatusDead || s == StatusLost
}

// ── Species ───────────────────────────────────────────────────────────────────

type Species string

const (
	SpeciesBovine  Species = "bovino"
	SpeciesOvine   Species = "ovino"
	SpeciesCaprine Species = "caprino"
	SpeciesPorcine Species = "porcino"
	SpeciesEquine  Species = "equino"
)

func (s Species) IsValid() bool {
	switch s {
	case SpeciesBovine, SpeciesOvine, SpeciesCaprine, SpeciesPorcine, SpeciesEquine:
		return true
	}
	return false
}

// ── Animal Entity ─────────────────────────────────────────────────────────────

type Animal struct {
	ID          uuid.UUID    `json:"id"`
	TenantID    uuid.UUID    `json:"tenant_id"`
	FarmID      uuid.UUID    `json:"farm_id"`
	VisualTag   string       `json:"visual_tag"`
	RfidTag     *string      `json:"rfid_tag,omitempty"`
	Species     Species      `json:"species"`
	Breed       *string      `json:"breed,omitempty"`
	Sex         string       `json:"sex"` // "M" | "F"
	BirthDate   time.Time    `json:"birth_date"`
	WeightKg    *float64     `json:"weight_kg,omitempty"`
	Color       *string      `json:"color,omitempty"`
	Status      AnimalStatus `json:"status"`
	SireID      *uuid.UUID   `json:"sire_id,omitempty"`
	DamID       *uuid.UUID   `json:"dam_id,omitempty"`
	PhotoURL    *string      `json:"photo_url,omitempty"`
	Notes       *string      `json:"notes,omitempty"`
	StatusDate  *time.Time   `json:"status_date,omitempty"`
	StatusNotes *string      `json:"status_notes,omitempty"`
	SyncVersion int64        `json:"sync_version"`
	IsDeleted   bool         `json:"is_deleted"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ── Pedigree Tree ─────────────────────────────────────────────────────────────

type PedigreeNode struct {
	Animal *Animal       `json:"animal"`
	Sire   *PedigreeNode `json:"sire"`
	Dam    *PedigreeNode `json:"dam"`
}
