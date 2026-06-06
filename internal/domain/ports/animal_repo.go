package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"mifinca-backend/internal/domain"
)

type ListAnimalsParams struct {
	TenantID uuid.UUID
	FarmID   *uuid.UUID
	Species  *string
	Sex      *string
	Status   *string
	Search   *string
	Limit    int32
	Cursor   *string // visual_tag del último para keyset pagination
}

type AnimalRepository interface {
	Create(ctx context.Context, a *domain.Animal) error
	CreateTx(ctx context.Context, tx pgx.Tx, a *domain.Animal) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Animal, error)
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) (*domain.Animal, error)
	Update(ctx context.Context, a *domain.Animal) error
	UpdateTx(ctx context.Context, tx pgx.Tx, a *domain.Animal) error
	PatchTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, fields map[string]any) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	SoftDeleteTx(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) error
	List(ctx context.Context, p ListAnimalsParams) ([]*domain.Animal, int64, error)
	GetPedigreeRows(ctx context.Context, tenantID, id uuid.UUID) ([]*domain.Animal, error)
	GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.Animal, error)
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
