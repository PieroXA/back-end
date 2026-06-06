package ports

import (
	"context"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
)

type ListHealthParams struct {
	TenantID  uuid.UUID
	AnimalID  *uuid.UUID
	EventType *string
	DateFrom  *string
	DateTo    *string
	Limit     int32
	Offset    int32
}

type HealthRepository interface {
	Create(ctx context.Context, e *domain.HealthEvent) error
	BulkCreate(ctx context.Context, events []*domain.HealthEvent) (int, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.HealthEvent, error)
	List(ctx context.Context, p ListHealthParams) ([]*domain.HealthEvent, error)
	GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.HealthEvent, error)
	GetUpcomingByNextDate(ctx context.Context, tenantID uuid.UUID, daysAhead int) ([]*domain.HealthEvent, error)
}
