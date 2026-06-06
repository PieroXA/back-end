package health

import (
	"time"

	"github.com/google/uuid"
)

type CreateHealthDTO struct {
	TenantID    uuid.UUID
	AnimalID    uuid.UUID
	EventType   string
	EventDate   time.Time
	ProductName *string
	DoseMl      *float64
	LotNumber   *string
	VetName     *string
	NextDate    *time.Time
	Notes       *string
	Cost        *float64
	BulkBatchID *uuid.UUID
}
