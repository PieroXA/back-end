package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
	infraAuth "mifinca-backend/internal/infrastructure/auth"
)

type Service struct {
	userRepo   ports.UserRepository
	tenantRepo ports.TenantRepository
	deviceRepo ports.DeviceRepository
	farmRepo   ports.FarmRepository
	jwtMaker   *infraAuth.JWTMaker
}

func NewService(
	userRepo ports.UserRepository,
	tenantRepo ports.TenantRepository,
	deviceRepo ports.DeviceRepository,
	farmRepo ports.FarmRepository,
	jwtMaker *infraAuth.JWTMaker,
) *Service {
	return &Service{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		deviceRepo: deviceRepo,
		farmRepo:   farmRepo,
		jwtMaker:   jwtMaker,
	}
}

// Register creates a new tenant, owner user, and default farm. Returns auth tokens.
func (s *Service) Register(ctx context.Context, dto RegisterDTO) (*AuthResult, error) {
	// Check email uniqueness
	existing, _ := s.userRepo.GetByEmail(ctx, dto.Email)
	if existing != nil {
		return nil, fmt.Errorf("auth.Register: email already in use: %w", domain.ErrConflict)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("auth.Register: hash password: %w", err)
	}

	now := time.Now()
	tenantID := uuid.New()
	userID := uuid.New()
	farmID := uuid.New()

	// Create tenant
	tenant := &domain.Tenant{
		ID:                tenantID,
		Name:              dto.FullName,
		Plan:              "basic",
		GlobalSyncVersion: 0,
		IsActive:          true,
		CreatedAt:         now,
	}
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("auth.Register: create tenant: %w", err)
	}

	// Create user (owner)
	user := &domain.User{
		ID:           userID,
		TenantID:     tenantID,
		Email:        dto.Email,
		PasswordHash: string(hash),
		FullName:     dto.FullName,
		Role:         domain.RoleOwner,
		IsActive:     true,
		CreatedAt:    now,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("auth.Register: create user: %w", err)
	}

	// Create default farm
	farm := &domain.Farm{
		ID:        farmID,
		TenantID:  tenantID,
		Name:      dto.FarmName,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.farmRepo.Create(ctx, farm); err != nil {
		return nil, fmt.Errorf("auth.Register: create farm: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtMaker.CreateAccessToken(userID, tenantID, string(domain.RoleOwner), "")
	if err != nil {
		return nil, fmt.Errorf("auth.Register: create access token: %w", err)
	}
	refreshToken, err := s.jwtMaker.CreateRefreshToken(userID, tenantID, string(domain.RoleOwner), "")
	if err != nil {
		return nil, fmt.Errorf("auth.Register: create refresh token: %w", err)
	}

	return &AuthResult{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserID:          userID.String(),
		TenantID:        tenantID.String(),
		Role:            string(domain.RoleOwner),
		LastSyncVersion: 0,
	}, nil
}

// Login authenticates a user and returns tokens.
func (s *Service) Login(ctx context.Context, dto LoginDTO) (*AuthResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, dto.Email)
	if err != nil {
		return nil, fmt.Errorf("auth.Login: %w", domain.ErrUnauthorized)
	}
	if !user.IsActive {
		return nil, fmt.Errorf("auth.Login: account deactivated: %w", domain.ErrForbidden)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, fmt.Errorf("auth.Login: invalid credentials: %w", domain.ErrUnauthorized)
	}

	// Upsert device
	now := time.Now()
	device := &domain.Device{
		ID:        uuid.New(),
		TenantID:  user.TenantID,
		UserID:    user.ID,
		DeviceUID: dto.DeviceID,
		LastSeenAt: &now,
		CreatedAt: now,
	}
	if err := s.deviceRepo.Upsert(ctx, device); err != nil {
		return nil, fmt.Errorf("auth.Login: upsert device: %w", err)
	}

	// Get last sync version for this device
	existingDevice, _ := s.deviceRepo.GetByDeviceUID(ctx, user.TenantID, dto.DeviceID)
	var lastSync int64
	if existingDevice != nil {
		lastSync = existingDevice.LastSyncVersion
	}

	// Generate tokens
	accessToken, err := s.jwtMaker.CreateAccessToken(user.ID, user.TenantID, string(user.Role), dto.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("auth.Login: create access token: %w", err)
	}
	refreshToken, err := s.jwtMaker.CreateRefreshToken(user.ID, user.TenantID, string(user.Role), dto.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("auth.Login: create refresh token: %w", err)
	}

	return &AuthResult{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserID:          user.ID.String(),
		TenantID:        user.TenantID.String(),
		Role:            string(user.Role),
		LastSyncVersion: lastSync,
	}, nil
}

// Refresh generates a new token pair from a valid refresh token.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	claims, err := s.jwtMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: %w", domain.ErrUnauthorized)
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || !user.IsActive {
		return nil, fmt.Errorf("auth.Refresh: %w", domain.ErrUnauthorized)
	}

	newAccess, err := s.jwtMaker.CreateAccessToken(user.ID, user.TenantID, string(user.Role), claims.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: create access token: %w", err)
	}
	newRefresh, err := s.jwtMaker.CreateRefreshToken(user.ID, user.TenantID, string(user.Role), claims.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: create refresh token: %w", err)
	}

	return &AuthResult{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		UserID:       user.ID.String(),
		TenantID:     user.TenantID.String(),
		Role:         string(user.Role),
	}, nil
}
