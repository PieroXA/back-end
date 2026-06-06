package movement

import (
	"time"

	"github.com/google/uuid"
)

type CreateMovementDTO struct {
	TenantID     uuid.UUID
	AnimalID     uuid.UUID
	FromFarmID   *uuid.UUID
	ToFarmID     *uuid.UUID
	MovementDate time.Time
	Reason       string
	Notes        *string
}
