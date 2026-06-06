package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"mifinca-backend/internal/domain"
)

type DeviceRepository struct {
	pool *pgxpool.Pool
}

func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{pool: pool}
}

func (r *DeviceRepository) Upsert(ctx context.Context, d *domain.Device) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO devices (id, tenant_id, user_id, device_uid, last_sync_version, last_seen_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (tenant_id, device_uid)
		 DO UPDATE SET user_id = $3, last_seen_at = $6`,
		d.ID, d.TenantID, d.UserID, d.DeviceUID, d.LastSyncVersion, d.LastSeenAt, d.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("deviceRepo.Upsert: %w", err)
	}
	return nil
}

func (r *DeviceRepository) GetByDeviceUID(ctx context.Context, tenantID uuid.UUID, deviceUID string) (*domain.Device, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, device_uid, last_sync_version, last_seen_at, created_at
		 FROM devices WHERE tenant_id = $1 AND device_uid = $2`,
		tenantID, deviceUID,
	)

	d := &domain.Device{}
	var lastSeen *time.Time
	err := row.Scan(&d.ID, &d.TenantID, &d.UserID, &d.DeviceUID, &d.LastSyncVersion, &lastSeen, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("deviceRepo.GetByDeviceUID: %w", domain.ErrNotFound)
	}
	d.LastSeenAt = lastSeen
	return d, nil
}

func (r *DeviceRepository) UpdateLastSyncVersion(ctx context.Context, tenantID uuid.UUID, deviceUID string, version int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE devices SET last_sync_version = $3, last_seen_at = NOW()
		 WHERE tenant_id = $1 AND device_uid = $2`,
		tenantID, deviceUID, version,
	)
	if err != nil {
		return fmt.Errorf("deviceRepo.UpdateLastSyncVersion: %w", err)
	}
	return nil
}
