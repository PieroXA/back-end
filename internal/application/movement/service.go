package movement

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
)

type Service struct {
	repo ports.MovementRepository
	bus  ports.EventBus
}

func NewService(repo ports.MovementRepository, bus ports.EventBus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
	}
}

func (s *Service) Create(ctx context.Context, dto CreateMovementDTO) (*domain.Movement, error) {
	reason := domain.MovementReason(dto.Reason)
	if !reason.IsValid() {
		return nil, fmt.Errorf("movement.Create: invalid reason: %w", domain.ErrValidation)
	}

	now := time.Now()
	movement := &domain.Movement{
		ID:           uuid.New(),
		TenantID:     dto.TenantID,
		AnimalID:     dto.AnimalID,
		FromFarmID:   dto.FromFarmID,
		ToFarmID:     dto.ToFarmID,
		MovementDate: dto.MovementDate,
		Reason:       reason,
		Notes:        dto.Notes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("movement.Create: repo failed: %w", err)
	}

	s.bus.Publish(ports.Event{
		EventType:  "MOVEMENT_CREATED",
		EntityType: "movement",
		EntityID:   movement.ID.String(),
		TenantID:   movement.TenantID.String(),
		Timestamp:  now.Format(time.RFC3339),
	})

	return movement, nil
}
