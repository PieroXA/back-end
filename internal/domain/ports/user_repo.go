package ports

import (
	"context"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type TenantRepository interface {
	Create(ctx context.Context, t *domain.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	IncrementSyncVersion(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

type DeviceRepository interface {
	Upsert(ctx context.Context, d *domain.Device) error
	GetByDeviceUID(ctx context.Context, tenantID uuid.UUID, deviceUID string) (*domain.Device, error)
	UpdateLastSyncVersion(ctx context.Context, tenantID uuid.UUID, deviceUID string, version int64) error
}
