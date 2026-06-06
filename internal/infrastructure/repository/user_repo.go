package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"mifinca-backend/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, tenant_id, email, password_hash, full_name, role, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.TenantID, u.Email, u.PasswordHash, u.FullName, u.Role, u.IsActive, u.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, full_name, role, is_active, created_at
		 FROM users WHERE email = $1`, email)

	u := &domain.User{}
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", domain.ErrNotFound)
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, full_name, role, is_active, created_at
		 FROM users WHERE id = $1`, id)

	u := &domain.User{}
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetByID: %w", domain.ErrNotFound)
	}
	return u, nil
}
