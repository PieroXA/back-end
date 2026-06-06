package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
	"mifinca-backend/internal/infrastructure/sqlc"
)

type MovementRepository struct {
	pool *pgxpool.Pool
}

func NewMovementRepository(pool *pgxpool.Pool) *MovementRepository {
	return &MovementRepository{pool: pool}
}

func (r *MovementRepository) withQuerier(tx pgx.Tx) sqlc.Querier {
	if tx != nil {
		return sqlc.New(tx)
	}
	return sqlc.New(r.pool)
}

func (r *MovementRepository) Create(ctx context.Context, m *domain.Movement) error {
	return r.CreateTx(ctx, nil, m)
}

func (r *MovementRepository) CreateTx(ctx context.Context, tx pgx.Tx, m *domain.Movement) error {
	q := r.withQuerier(tx)
	
	params := sqlc.CreateMovementParams{
		ID:           pgtype.UUID{Bytes: m.ID, Valid: true},
		TenantID:     pgtype.UUID{Bytes: m.TenantID, Valid: true},
		AnimalID:     pgtype.UUID{Bytes: m.AnimalID, Valid: true},
		MovementDate: pgtype.Date{Time: m.MovementDate, Valid: true},
		Reason:       string(m.Reason),
	}
	if m.FromFarmID != nil {
		params.FromFarmID = pgtype.UUID{Bytes: *m.FromFarmID, Valid: true}
	}
	if m.ToFarmID != nil {
		params.ToFarmID = pgtype.UUID{Bytes: *m.ToFarmID, Valid: true}
	}
	if m.Notes != nil {
		params.Notes = pgtype.Text{String: *m.Notes, Valid: true}
	}

	row, err := q.CreateMovement(ctx, params)
	if err != nil {
		return fmt.Errorf("movementRepo.CreateTx: %w", err)
	}

	m.SyncVersion = row.SyncVersion
	m.CreatedAt = row.CreatedAt.Time
	m.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *MovementRepository) GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.Movement, error) {
	q := r.withQuerier(nil)
	rows, err := q.GetMovementsSince(ctx, sqlc.GetMovementsSinceParams{
		TenantID:    pgtype.UUID{Bytes: tenantID, Valid: true},
		SyncVersion: since,
	})
	if err != nil {
		return nil, fmt.Errorf("movementRepo.GetSinceVersion: %w", err)
	}

	var movements []*domain.Movement
	for _, row := range rows {
		movements = append(movements, mapSqlcMovementToDomain(row))
	}
	return movements, nil
}

func mapSqlcMovementToDomain(row sqlc.Movement) *domain.Movement {
	m := &domain.Movement{
		ID:           row.ID.Bytes,
		TenantID:     row.TenantID.Bytes,
		AnimalID:     row.AnimalID.Bytes,
		MovementDate: row.MovementDate.Time,
		Reason:       domain.MovementReason(row.Reason),
		SyncVersion:  row.SyncVersion,
		IsDeleted:    row.IsDeleted,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
	if row.FromFarmID.Valid {
		id := uuid.UUID(row.FromFarmID.Bytes)
		m.FromFarmID = &id
	}
	if row.ToFarmID.Valid {
		id := uuid.UUID(row.ToFarmID.Bytes)
		m.ToFarmID = &id
	}
	if row.Notes.Valid {
		m.Notes = &row.Notes.String
	}
	return m
}

// Stub for remaining ports (UC-04.2)
func (r *MovementRepository) List(ctx context.Context, p ports.ListMovementsParams) ([]*domain.Movement, error) {
	return nil, fmt.Errorf("not implemented")
}
