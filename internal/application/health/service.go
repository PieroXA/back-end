package health

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
)

type Service struct {
	repo ports.HealthRepository
	bus  ports.EventBus
}

func NewService(repo ports.HealthRepository, bus ports.EventBus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
	}
}

func (s *Service) Create(ctx context.Context, dto CreateHealthDTO) (*domain.HealthEvent, error) {
	evtType := domain.EventType(dto.EventType)
	if !evtType.IsValid() {
		return nil, fmt.Errorf("health.Create: invalid event_type: %w", domain.ErrValidation)
	}
	if evtType.RequiresProduct() && dto.ProductName == nil {
		return nil, fmt.Errorf("health.Create: event_type requires product_name: %w", domain.ErrValidation)
	}

	now := time.Now()
	event := &domain.HealthEvent{
		ID:          uuid.New(),
		TenantID:    dto.TenantID,
		AnimalID:    dto.AnimalID,
		EventType:   evtType,
		EventDate:   dto.EventDate,
		ProductName: dto.ProductName,
		DoseMl:      dto.DoseMl,
		LotNumber:   dto.LotNumber,
		VetName:     dto.VetName,
		NextDate:    dto.NextDate,
		Notes:       dto.Notes,
		Cost:        dto.Cost,
		BulkBatchID: dto.BulkBatchID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("health.Create: repo failed: %w", err)
	}

	s.bus.Publish(ports.Event{
		EventType:  "HEALTH_CREATED",
		EntityType: "health_event",
		EntityID:   event.ID.String(),
		TenantID:   event.TenantID.String(),
		Timestamp:  now.Format(time.RFC3339),
	})

	return event, nil
}
