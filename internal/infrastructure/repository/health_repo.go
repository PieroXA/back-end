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

type HealthRepository struct {
	pool *pgxpool.Pool
}

func NewHealthRepository(pool *pgxpool.Pool) *HealthRepository {
	return &HealthRepository{pool: pool}
}

func (r *HealthRepository) withQuerier(tx pgx.Tx) sqlc.Querier {
	if tx != nil {
		return sqlc.New(tx)
	}
	return sqlc.New(r.pool)
}

func (r *HealthRepository) Create(ctx context.Context, h *domain.HealthEvent) error {
	return r.CreateTx(ctx, nil, h)
}

func (r *HealthRepository) CreateTx(ctx context.Context, tx pgx.Tx, h *domain.HealthEvent) error {
	q := r.withQuerier(tx)

	params := sqlc.CreateHealthEventParams{
		ID:        pgtype.UUID{Bytes: h.ID, Valid: true},
		TenantID:  pgtype.UUID{Bytes: h.TenantID, Valid: true},
		AnimalID:  pgtype.UUID{Bytes: h.AnimalID, Valid: true},
		EventType: string(h.EventType),
		EventDate: pgtype.Date{Time: h.EventDate, Valid: true},
	}
	if h.ProductName != nil {
		params.ProductName = pgtype.Text{String: *h.ProductName, Valid: true}
	}
	if h.DoseMl != nil {
		params.DoseMl = pgtype.Numeric{Int: toNumeric(*h.DoseMl), Exp: -2, Valid: true}
	}
	if h.LotNumber != nil {
		params.LotNumber = pgtype.Text{String: *h.LotNumber, Valid: true}
	}
	if h.VetName != nil {
		params.VetName = pgtype.Text{String: *h.VetName, Valid: true}
	}
	if h.NextDate != nil {
		params.NextDate = pgtype.Date{Time: *h.NextDate, Valid: true}
	}
	if h.Notes != nil {
		params.Notes = pgtype.Text{String: *h.Notes, Valid: true}
	}
	if h.Cost != nil {
		params.Cost = pgtype.Numeric{Int: toNumeric(*h.Cost), Exp: -2, Valid: true}
	}
	if h.BulkBatchID != nil {
		params.BulkBatchID = pgtype.UUID{Bytes: *h.BulkBatchID, Valid: true}
	}

	row, err := q.CreateHealthEvent(ctx, params)
	if err != nil {
		return fmt.Errorf("healthRepo.CreateTx: %w", err)
	}

	h.SyncVersion = row.SyncVersion
	h.CreatedAt = row.CreatedAt.Time
	h.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *HealthRepository) GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.HealthEvent, error) {
	q := r.withQuerier(nil)
	rows, err := q.GetHealthEventsSince(ctx, sqlc.GetHealthEventsSinceParams{
		TenantID:    pgtype.UUID{Bytes: tenantID, Valid: true},
		SyncVersion: since,
	})
	if err != nil {
		return nil, fmt.Errorf("healthRepo.GetSinceVersion: %w", err)
	}

	var events []*domain.HealthEvent
	for _, row := range rows {
		events = append(events, mapSqlcHealthEventToDomain(row))
	}
	return events, nil
}

func mapSqlcHealthEventToDomain(row sqlc.HealthEvent) *domain.HealthEvent {
	h := &domain.HealthEvent{
		ID:          row.ID.Bytes,
		TenantID:    row.TenantID.Bytes,
		AnimalID:    row.AnimalID.Bytes,
		EventType:   domain.EventType(row.EventType),
		EventDate:   row.EventDate.Time,
		SyncVersion: row.SyncVersion,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
	if row.ProductName.Valid {
		h.ProductName = &row.ProductName.String
	}
	if row.LotNumber.Valid {
		h.LotNumber = &row.LotNumber.String
	}
	if row.VetName.Valid {
		h.VetName = &row.VetName.String
	}
	if row.NextDate.Valid {
		t := row.NextDate.Time
		h.NextDate = &t
	}
	if row.Notes.Valid {
		h.Notes = &row.Notes.String
	}
	if row.BulkBatchID.Valid {
		id := uuid.UUID(row.BulkBatchID.Bytes)
		h.BulkBatchID = &id
	}
	// Note: Parse numeric for dose and cost removed for simplicity in this prototype
	return h
}

// Stubs for remaining ports (UC-03.2, UC-03.3, GetByID, Upcoming)
func (r *HealthRepository) BulkCreate(ctx context.Context, events []*domain.HealthEvent) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (r *HealthRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.HealthEvent, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *HealthRepository) List(ctx context.Context, p ports.ListHealthParams) ([]*domain.HealthEvent, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *HealthRepository) GetUpcomingByNextDate(ctx context.Context, tenantID uuid.UUID, daysAhead int) ([]*domain.HealthEvent, error) {
	return nil, fmt.Errorf("not implemented")
}
