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

func (r *UserRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, email, password_hash, full_name, role, is_active, created_at
		 FROM users WHERE tenant_id = $1 ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("userRepo.List: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("userRepo.List.Scan: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	commandTag, err := r.pool.Exec(ctx,
		`UPDATE users 
		 SET email = $1, password_hash = $2, full_name = $3, role = $4, is_active = $5
		 WHERE id = $6 AND tenant_id = $7`,
		u.Email, u.PasswordHash, u.FullName, u.Role, u.IsActive, u.ID, u.TenantID,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Update: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("userRepo.Update: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error {
	commandTag, err := r.pool.Exec(ctx,
		`UPDATE users SET is_active = false WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Delete: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("userRepo.Delete: %w", domain.ErrNotFound)
	}
	return nil
}
