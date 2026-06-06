package ports

import (
	"context"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
)

type FarmRepository interface {
	Create(ctx context.Context, f *domain.Farm) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Farm, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Farm, error)
	Update(ctx context.Context, f *domain.Farm) error
	GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.Farm, error)
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
