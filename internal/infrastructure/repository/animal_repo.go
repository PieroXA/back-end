package repository

import (
	"context"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
	"mifinca-backend/internal/infrastructure/sqlc"
)

type AnimalRepository struct {
	pool *pgxpool.Pool
}

func NewAnimalRepository(pool *pgxpool.Pool) *AnimalRepository {
	return &AnimalRepository{pool: pool}
}

// withQuerier allows choosing between the pool or a transaction
func (r *AnimalRepository) withQuerier(tx pgx.Tx) sqlc.Querier {
	if tx != nil {
		return sqlc.New(tx)
	}
	return sqlc.New(r.pool)
}

func (r *AnimalRepository) Create(ctx context.Context, a *domain.Animal) error {
	return r.CreateTx(ctx, nil, a)
}

func (r *AnimalRepository) CreateTx(ctx context.Context, tx pgx.Tx, a *domain.Animal) error {
	q := r.withQuerier(tx)
	params := sqlc.CreateAnimalParams{
		ID:        pgtype.UUID{Bytes: a.ID, Valid: true},
		TenantID:  pgtype.UUID{Bytes: a.TenantID, Valid: true},
		FarmID:    pgtype.UUID{Bytes: a.FarmID, Valid: true},
		VisualTag: a.VisualTag,
		Species:   string(a.Species),
		Sex:       a.Sex,
		BirthDate: pgtype.Date{Time: a.BirthDate, Valid: true},
	}
	if a.RfidTag != nil {
		params.RfidTag = pgtype.Text{String: *a.RfidTag, Valid: true}
	}
	if a.Breed != nil {
		params.Breed = pgtype.Text{String: *a.Breed, Valid: true}
	}
	if a.WeightKg != nil {
		params.WeightKg = pgtype.Numeric{Int: toNumeric(*a.WeightKg), Exp: -2, Valid: true} // Simplified for prototype
	}
	if a.Color != nil {
		params.Color = pgtype.Text{String: *a.Color, Valid: true}
	}
	if a.SireID != nil {
		params.SireID = pgtype.UUID{Bytes: *a.SireID, Valid: true}
	}
	if a.DamID != nil {
		params.DamID = pgtype.UUID{Bytes: *a.DamID, Valid: true}
	}
	if a.Notes != nil {
		params.Notes = pgtype.Text{String: *a.Notes, Valid: true}
	}

	row, err := q.CreateAnimal(ctx, params)
	if err != nil {
		return fmt.Errorf("animalRepo.CreateTx: %w", err)
	}

	a.SyncVersion = row.SyncVersion
	a.CreatedAt = row.CreatedAt.Time
	a.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *AnimalRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Animal, error) {
	q := r.withQuerier(nil)
	row, err := q.GetAnimalByID(ctx, sqlc.GetAnimalByIDParams{
		TenantID: pgtype.UUID{Bytes: tenantID, Valid: true},
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("animalRepo.GetByID: %w", domain.ErrNotFound)
	}
	return mapSqlcAnimalToDomain(row), nil
}

func (r *AnimalRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) (*domain.Animal, error) {
	q := r.withQuerier(tx)
	row, err := q.GetAnimalByIDForUpdate(ctx, sqlc.GetAnimalByIDForUpdateParams{
		TenantID: pgtype.UUID{Bytes: tenantID, Valid: true},
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("animalRepo.GetByIDForUpdate: %w", domain.ErrNotFound)
	}
	return mapSqlcAnimalToDomain(row), nil
}

func (r *AnimalRepository) Update(ctx context.Context, a *domain.Animal) error {
	return r.UpdateTx(ctx, nil, a)
}

func (r *AnimalRepository) UpdateTx(ctx context.Context, tx pgx.Tx, a *domain.Animal) error {
	// A full update function matching domain to DB schema.
	// Since sqlc didn't generate a generic Update (we skipped writing the massive query), 
	// we use a pgx exec manually to not over-complicate the sqlc generator.
	conn := r.pool
	var querier interface{ Exec(context.Context, string, ...any) (pgconn.CommandTag, error) } = conn
	if tx != nil {
		querier = tx
	}

	_, err := querier.Exec(ctx,
		`UPDATE animals SET 
			farm_id=$3, visual_tag=$4, rfid_tag=$5, species=$6, breed=$7, sex=$8, 
			birth_date=$9, weight_kg=$10, color=$11, status=$12, sire_id=$13, dam_id=$14, 
			photo_url=$15, notes=$16, status_date=$17, status_notes=$18, updated_at=NOW()
		 WHERE id=$1 AND tenant_id=$2 AND NOT is_deleted`,
		a.ID, a.TenantID, a.FarmID, a.VisualTag, a.RfidTag, a.Species, a.Breed, a.Sex,
		a.BirthDate, a.WeightKg, a.Color, a.Status, a.SireID, a.DamID,
		a.PhotoURL, a.Notes, a.StatusDate, a.StatusNotes,
	)
	if err != nil {
		return fmt.Errorf("animalRepo.UpdateTx: %w", err)
	}
	return nil
}

func (r *AnimalRepository) PatchTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, fields map[string]any) error {
	// Implementation for Sync Engine patch (LWW), to be detailed in Sprint 3
	return nil
}

func (r *AnimalRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.SoftDeleteTx(ctx, nil, tenantID, id)
}

func (r *AnimalRepository) SoftDeleteTx(ctx context.Context, tx pgx.Tx, tenantID, id uuid.UUID) error {
	q := r.withQuerier(tx)
	err := q.SoftDeleteAnimal(ctx, sqlc.SoftDeleteAnimalParams{
		TenantID: pgtype.UUID{Bytes: tenantID, Valid: true},
		ID:       pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("animalRepo.SoftDeleteTx: %w", err)
	}
	return nil
}

func (r *AnimalRepository) List(ctx context.Context, p ports.ListAnimalsParams) ([]*domain.Animal, int64, error) {
	q := r.withQuerier(nil)
	
	params := sqlc.ListAnimalsParams{
		TenantID: pgtype.UUID{Bytes: p.TenantID, Valid: true},
		Limit:    p.Limit,
	}
	if p.FarmID != nil {
		params.FarmID = pgtype.UUID{Bytes: *p.FarmID, Valid: true}
	}
	if p.Species != nil {
		params.Species = pgtype.Text{String: *p.Species, Valid: true}
	}
	if p.Sex != nil {
		params.Sex = pgtype.Text{String: *p.Sex, Valid: true}
	}
	if p.Status != nil {
		params.Status = pgtype.Text{String: *p.Status, Valid: true}
	}
	if p.Search != nil {
		params.Search = pgtype.Text{String: *p.Search, Valid: true}
	}
	if p.Cursor != nil {
		params.Cursor = pgtype.Text{String: *p.Cursor, Valid: true}
	}

	rows, err := q.ListAnimals(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("animalRepo.List: %w", err)
	}

	var animals []*domain.Animal
	for _, row := range rows {
		animals = append(animals, mapSqlcAnimalToDomain(row))
	}

	count, err := q.CountAnimalsByTenant(ctx, pgtype.UUID{Bytes: p.TenantID, Valid: true})
	if err != nil {
		return nil, 0, fmt.Errorf("animalRepo.List count: %w", err)
	}

	return animals, count, nil
}

func (r *AnimalRepository) GetPedigreeRows(ctx context.Context, tenantID, id uuid.UUID) ([]*domain.Animal, error) {
	q := r.withQuerier(nil)
	rows, err := q.GetPedigreeRows(ctx, sqlc.GetPedigreeRowsParams{
		ID:       pgtype.UUID{Bytes: id, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("animalRepo.GetPedigreeRows: %w", err)
	}

	var animals []*domain.Animal
	for _, row := range rows {
		a := &domain.Animal{
			ID:        row.ID.Bytes,
			VisualTag: row.VisualTag,
			Species:   domain.Species(row.Species),
			Sex:       row.Sex,
			BirthDate: row.BirthDate.Time,
		}
		if row.Breed.Valid {
			a.Breed = &row.Breed.String
		}
		if row.SireID.Valid {
			id := uuid.UUID(row.SireID.Bytes)
			a.SireID = &id
		}
		if row.DamID.Valid {
			id := uuid.UUID(row.DamID.Bytes)
			a.DamID = &id
		}
		animals = append(animals, a)
	}
	return animals, nil
}

func (r *AnimalRepository) GetSinceVersion(ctx context.Context, tenantID uuid.UUID, since int64) ([]*domain.Animal, error) {
	q := r.withQuerier(nil)
	rows, err := q.GetAnimalsSince(ctx, sqlc.GetAnimalsSinceParams{
		TenantID:    pgtype.UUID{Bytes: tenantID, Valid: true},
		SyncVersion: since,
	})
	if err != nil {
		return nil, fmt.Errorf("animalRepo.GetSinceVersion: %w", err)
	}

	var animals []*domain.Animal
	for _, row := range rows {
		animals = append(animals, mapSqlcAnimalToDomain(row))
	}
	return animals, nil
}

func (r *AnimalRepository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	q := r.withQuerier(nil)
	count, err := q.CountAnimalsByTenant(ctx, pgtype.UUID{Bytes: tenantID, Valid: true})
	return count, err
}

func toNumeric(f float64) *big.Int {
	val := new(big.Float).SetFloat64(f * 100)
	i, _ := val.Int(nil)
	return i
}

func mapSqlcAnimalToDomain(row sqlc.Animal) *domain.Animal {
	a := &domain.Animal{
		ID:          row.ID.Bytes,
		TenantID:    row.TenantID.Bytes,
		FarmID:      row.FarmID.Bytes,
		VisualTag:   row.VisualTag,
		Species:     domain.Species(row.Species),
		Sex:         row.Sex,
		BirthDate:   row.BirthDate.Time,
		Status:      domain.AnimalStatus(row.Status),
		SyncVersion: row.SyncVersion,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	if row.RfidTag.Valid {
		a.RfidTag = &row.RfidTag.String
	}
	if row.Breed.Valid {
		a.Breed = &row.Breed.String
	}
	if row.Color.Valid {
		a.Color = &row.Color.String
	}
	if row.SireID.Valid {
		id := uuid.UUID(row.SireID.Bytes)
		a.SireID = &id
	}
	if row.DamID.Valid {
		id := uuid.UUID(row.DamID.Bytes)
		a.DamID = &id
	}
	if row.Notes.Valid {
		a.Notes = &row.Notes.String
	}
	if row.StatusDate.Valid {
		t := row.StatusDate.Time
		a.StatusDate = &t
	}
	if row.StatusNotes.Valid {
		a.StatusNotes = &row.StatusNotes.String
	}

	return a
}
