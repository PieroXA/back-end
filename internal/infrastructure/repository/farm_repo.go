package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"mifinca-backend/internal/domain"
)

type FarmRepository struct {
	pool *pgxpool.Pool
}

func NewFarmRepository(pool *pgxpool.Pool) *FarmRepository {
	return &FarmRepository{pool: pool}
}

func (r *FarmRepository) Create(ctx context.Context, f *domain.Farm) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO farms (id, tenant_id, name, address, latitude, longitude, notes, sync_version, is_deleted, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		f.ID, f.TenantID, f.Name, f.Address, f.Latitude, f.Longitude, f.Notes,
		f.SyncVersion, f.IsDeleted, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("farmRepo.Create: %w", err)
	}
	return nil
}

func (r *FarmRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Farm, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, address, latitude, longitude, notes, sync_version, is_deleted, created_at, updated_at
		 FROM farms WHERE id = $1 AND tenant_id = $2 AND NOT is_deleted`, id, tenantID)

	f := &domain.Farm{}
	err := row.Scan(&f.ID, &f.TenantID, &f.Name, &f.Address, &f.Latitude, &f.Longitude, &f.Notes,
		&f.SyncVersion, &f.IsDeleted, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("farmRepo.GetByID: %w", domain.ErrNotFound)
	}
	return f, nil
}

func (r *FarmRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.Farm, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, name, address, latitude, longitude, notes, sync_version, is_deleted, created_at, updated_at
		 FROM farms WHERE tenant_id = $1 AND NOT is_deleted ORDER BY name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("farmRepo.ListByTenant: %w", err)
	}
	defer rows.Close()

	var farms []*domain.Farm
	for rows.Next() {
		f := &domain.Farm{}
		if err := rows.Scan(&f.ID, &f.TenantID, &f.Name, &f.Address, &f.Latitude, &f.Longitude, &f.Notes,
			&f.SyncVersion, &f.IsDeleted, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("farmRepo.ListByTenant scan: %w", err)
		}
		farms = append(farms, f)
	}
	return farms, nil
}

func (r *FarmRepository) Update(ctx context.Context, f *domain.Farm) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE farms SET name=$2, address=$3, latitude=$4, longitude=$5, notes=$6, updated_at=NOW()
		 WHERE id = $1 AND tenant_id = $7`,
		f.ID, f.Name, f.Address, f.Latitude, f.Longitude, f.Notes, f.TenantID,
	)
	if err != nil {
		return fmt.Errorf("farmRepo.Update: %w", err)
	}
	return nil
}

func (r *FarmRepository) GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.Farm, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, name, address, latitude, longitude, notes, sync_version, is_deleted, created_at, updated_at
		 FROM farms WHERE tenant_id = $1 AND sync_version > $2 ORDER BY sync_version ASC`,
		tenantID, since)
	if err != nil {
		return nil, fmt.Errorf("farmRepo.GetSinceVersion: %w", err)
	}
	defer rows.Close()

	var farms []*domain.Farm
	for rows.Next() {
		f := &domain.Farm{}
		if err := rows.Scan(&f.ID, &f.TenantID, &f.Name, &f.Address, &f.Latitude, &f.Longitude, &f.Notes,
			&f.SyncVersion, &f.IsDeleted, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("farmRepo.GetSinceVersion scan: %w", err)
		}
		farms = append(farms, f)
	}
	return farms, nil
}

func (r *FarmRepository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM farms WHERE tenant_id = $1 AND NOT is_deleted`, tenantID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("farmRepo.CountByTenant: %w", err)
	}
	return count, nil
}
