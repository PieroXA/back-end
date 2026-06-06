package ports

import (
	"context"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
)

type ListMovementsParams struct {
	TenantID uuid.UUID
	AnimalID *uuid.UUID
	FarmID   *uuid.UUID
	DateFrom *string
	DateTo   *string
	Limit    int32
	Offset   int32
}

type MovementRepository interface {
	Create(ctx context.Context, m *domain.Movement) error
	List(ctx context.Context, p ListMovementsParams) ([]*domain.Movement, error)
	GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.Movement, error)
}
