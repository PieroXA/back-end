package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"mifinca-backend/internal/domain"
)

type TenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

func (r *TenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tenants (id, name, plan, global_sync_version, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.Name, t.Plan, t.GlobalSyncVersion, t.IsActive, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("tenantRepo.Create: %w", err)
	}
	return nil
}

func (r *TenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, plan, global_sync_version, is_active, created_at
		 FROM tenants WHERE id = $1`, id)

	t := &domain.Tenant{}
	err := row.Scan(&t.ID, &t.Name, &t.Plan, &t.GlobalSyncVersion, &t.IsActive, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("tenantRepo.GetByID: %w", domain.ErrNotFound)
	}
	return t, nil
}

func (r *TenantRepository) IncrementSyncVersion(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var version int64
	err := r.pool.QueryRow(ctx,
		`SELECT increment_tenant_sync_version($1)`, tenantID,
	).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("tenantRepo.IncrementSyncVersion: %w", err)
	}
	return version, nil
}
