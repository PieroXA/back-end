package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, tenantID uuid.UUID, req CreateUserRequest) (UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         domain.Role(req.Role),
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func (s *UserService) List(ctx context.Context, tenantID uuid.UUID) ([]UserResponse, error) {
	users, err := s.userRepo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var res []UserResponse
	for _, u := range users {
		res = append(res, toUserResponse(u))
	}
	return res, nil
}

func (s *UserService) GetByID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (UserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	if u.TenantID != tenantID {
		return UserResponse{}, domain.ErrNotFound
	}

	return toUserResponse(u), nil
}

func (s *UserService) Update(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req UpdateUserRequest) (UserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	if u.TenantID != tenantID {
		return UserResponse{}, domain.ErrNotFound
	}

	u.Email = req.Email
	u.FullName = req.FullName
	u.Role = domain.Role(req.Role)
	u.IsActive = *req.IsActive

	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return UserResponse{}, fmt.Errorf("failed to hash password: %w", err)
		}
		u.PasswordHash = string(hash)
	}

	if err := s.userRepo.Update(ctx, u); err != nil {
		return UserResponse{}, err
	}

	return toUserResponse(u), nil
}

func (s *UserService) Delete(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if u.TenantID != tenantID {
		return domain.ErrNotFound
	}

	return s.userRepo.Delete(ctx, id, tenantID)
}

func toUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
