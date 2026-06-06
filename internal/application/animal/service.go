package animal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mifinca-backend/internal/domain"
	"mifinca-backend/internal/domain/ports"
)

type Service struct {
	repo     ports.AnimalRepository
	farmRepo ports.FarmRepository
	moveRepo ports.MovementRepository
	bus      ports.EventBus
}

func NewService(
	repo ports.AnimalRepository,
	farmRepo ports.FarmRepository,
	moveRepo ports.MovementRepository,
	bus ports.EventBus,
) *Service {
	return &Service{
		repo:     repo,
		farmRepo: farmRepo,
		moveRepo: moveRepo,
		bus:      bus,
	}
}

// Create implements UC-02.1 Registrar Animal
func (s *Service) Create(ctx context.Context, dto CreateAnimalDTO) (*domain.Animal, error) {
	// Domain validation for species and sex
	species := domain.Species(dto.Species)
	if !species.IsValid() {
		return nil, fmt.Errorf("animal.Create: invalid species: %w", domain.ErrValidation)
	}
	if dto.Sex != "M" && dto.Sex != "F" {
		return nil, fmt.Errorf("animal.Create: invalid sex: %w", domain.ErrValidation)
	}

	// Verify farm exists and belongs to tenant
	farm, err := s.farmRepo.GetByID(ctx, dto.TenantID, dto.FarmID)
	if err != nil {
		return nil, fmt.Errorf("animal.Create: invalid farm: %w", err)
	}
	if farm == nil {
		return nil, fmt.Errorf("animal.Create: farm not found: %w", domain.ErrValidation)
	}

	// Sire & Dam validation (if provided)
	if dto.SireID != nil {
		sire, err := s.repo.GetByID(ctx, dto.TenantID, *dto.SireID)
		if err != nil || sire.Sex != "M" || sire.Species != species {
			return nil, fmt.Errorf("animal.Create: invalid sire: %w", domain.ErrValidation)
		}
	}
	if dto.DamID != nil {
		dam, err := s.repo.GetByID(ctx, dto.TenantID, *dto.DamID)
		if err != nil || dam.Sex != "F" || dam.Species != species {
			return nil, fmt.Errorf("animal.Create: invalid dam: %w", domain.ErrValidation)
		}
	}

	now := time.Now()
	animal := &domain.Animal{
		ID:        uuid.New(),
		TenantID:  dto.TenantID,
		FarmID:    dto.FarmID,
		VisualTag: dto.VisualTag,
		RfidTag:   dto.RfidTag,
		Species:   species,
		Breed:     dto.Breed,
		Sex:       dto.Sex,
		BirthDate: dto.BirthDate,
		WeightKg:  dto.WeightKg,
		Color:     dto.Color,
		Status:    domain.StatusActive,
		SireID:    dto.SireID,
		DamID:     dto.DamID,
		Notes:     dto.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, animal); err != nil {
		return nil, fmt.Errorf("animal.Create: repository failed: %w", err)
	}

	// Publish Event
	// Note: serialization of payload is omitted for brevity but should marshal the created entity
	s.bus.Publish(ports.Event{
		EventType:  "ANIMAL_CREATED",
		EntityType: "animal",
		EntityID:   animal.ID.String(),
		TenantID:   animal.TenantID.String(),
		Timestamp:  now.Format(time.RFC3339),
	})

	return animal, nil
}

// GetByID returns a single animal
func (s *Service) GetByID(ctx context.Context, tenantID, animalID uuid.UUID) (*domain.Animal, error) {
	return s.repo.GetByID(ctx, tenantID, animalID)
}

// List implements UC-02.5 Listar Animales
func (s *Service) List(ctx context.Context, p ports.ListAnimalsParams) ([]*domain.Animal, int64, error) {
	return s.repo.List(ctx, p)
}

// SoftDelete implements UC-02.4 Baja Lógica
func (s *Service) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	err := s.repo.SoftDelete(ctx, tenantID, id)
	if err != nil {
		return err
	}
	
	s.bus.Publish(ports.Event{
		EventType:  "ANIMAL_DELETED",
		EntityType: "animal",
		EntityID:   id.String(),
		TenantID:   tenantID.String(),
		Timestamp:  time.Now().Format(time.RFC3339),
	})
	return nil
}

// BuildPedigreeTree constructs the pedigree graph recursively from raw rows (UC-02.6 & UC-05)
func (s *Service) GetPedigree(ctx context.Context, tenantID, animalID uuid.UUID) (*domain.PedigreeNode, error) {
	rows, err := s.repo.GetPedigreeRows(ctx, tenantID, animalID)
	if err != nil {
		return nil, fmt.Errorf("animal.GetPedigree: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("animal.GetPedigree: root animal not found: %w", domain.ErrNotFound)
	}

	// Create a map to quickly look up animal objects by ID
	animalMap := make(map[uuid.UUID]*domain.Animal)
	for _, a := range rows {
		animalMap[a.ID] = a
	}

	// Helper recursive function to build the tree
	var build func(id uuid.UUID, depth int) *domain.PedigreeNode
	build = func(id uuid.UUID, depth int) *domain.PedigreeNode {
		if depth > 4 {
			return nil // Limit pedigree to 4 generations
		}
		a, ok := animalMap[id]
		if !ok {
			return nil
		}
		node := &domain.PedigreeNode{
			Animal: a,
		}
		if a.SireID != nil {
			node.Sire = build(*a.SireID, depth+1)
		}
		if a.DamID != nil {
			node.Dam = build(*a.DamID, depth+1)
		}
		return node
	}

	return build(animalID, 1), nil
}

// Update implements UC-02.2 Actualizar Animal (PATCH)
func (s *Service) Update(ctx context.Context, dto UpdateAnimalDTO) (*domain.Animal, error) {
	animal, err := s.repo.GetByID(ctx, dto.TenantID, dto.AnimalID)
	if err != nil {
		return nil, fmt.Errorf("animal.Update: not found: %w", err)
	}
	if animal.Status.IsFinal() {
		return nil, fmt.Errorf("animal.Update: cannot update an animal in a final status: %w", domain.ErrValidation)
	}

	if dto.VisualTag != nil {
		animal.VisualTag = *dto.VisualTag
	}
	if dto.RfidTag != nil {
		animal.RfidTag = dto.RfidTag
	}
	if dto.Breed != nil {
		animal.Breed = dto.Breed
	}
	if dto.WeightKg != nil {
		animal.WeightKg = dto.WeightKg
	}
	if dto.Color != nil {
		animal.Color = dto.Color
	}
	if dto.Notes != nil {
		animal.Notes = dto.Notes
	}

	// Validate Sires/Dams if they are being updated
	if dto.SireID != nil {
		sire, err := s.repo.GetByID(ctx, dto.TenantID, *dto.SireID)
		if err != nil || sire.Sex != "M" || sire.Species != animal.Species {
			return nil, fmt.Errorf("animal.Update: invalid sire: %w", domain.ErrValidation)
		}
		animal.SireID = dto.SireID
	}
	if dto.DamID != nil {
		dam, err := s.repo.GetByID(ctx, dto.TenantID, *dto.DamID)
		if err != nil || dam.Sex != "F" || dam.Species != animal.Species {
			return nil, fmt.Errorf("animal.Update: invalid dam: %w", domain.ErrValidation)
		}
		animal.DamID = dto.DamID
	}

	// Handle Farm change -> triggers implicit Movement
	var movement *domain.Movement
	now := time.Now()
	if dto.FarmID != nil && *dto.FarmID != animal.FarmID {
		oldFarm := animal.FarmID
		newFarm := *dto.FarmID
		
		farm, err := s.farmRepo.GetByID(ctx, dto.TenantID, newFarm)
		if err != nil || farm == nil {
			return nil, fmt.Errorf("animal.Update: invalid new farm: %w", domain.ErrValidation)
		}

		movement = &domain.Movement{
			ID:           uuid.New(),
			TenantID:     dto.TenantID,
			AnimalID:     animal.ID,
			FromFarmID:   &oldFarm,
			ToFarmID:     &newFarm,
			MovementDate: now,
			Reason:       domain.ReasonSaleTransfer,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		animal.FarmID = newFarm
	}

	animal.UpdatedAt = now

	// We use the basic Update here. In a real scenario with the Movement, we'd use UpdateTx inside a pgx.Tx
	// For Sprint 2 completeness without Tx Orchestrator yet, we do sequential calls.
	if err := s.repo.Update(ctx, animal); err != nil {
		return nil, fmt.Errorf("animal.Update: repository failed: %w", err)
	}

	if movement != nil {
		_ = s.moveRepo.Create(ctx, movement)
		s.bus.Publish(ports.Event{
			EventType:  "MOVEMENT_CREATED",
			EntityType: "movement",
			EntityID:   movement.ID.String(),
			TenantID:   movement.TenantID.String(),
			Timestamp:  now.Format(time.RFC3339),
		})
	}

	s.bus.Publish(ports.Event{
		EventType:  "ANIMAL_UPDATED",
		EntityType: "animal",
		EntityID:   animal.ID.String(),
		TenantID:   animal.TenantID.String(),
		Timestamp:  now.Format(time.RFC3339),
	})

	return animal, nil
}

// ChangeStatus implements UC-02.3 Cambiar Estado
func (s *Service) ChangeStatus(ctx context.Context, dto ChangeStatusDTO) (*domain.Animal, error) {
	animal, err := s.repo.GetByID(ctx, dto.TenantID, dto.AnimalID)
	if err != nil {
		return nil, fmt.Errorf("animal.ChangeStatus: not found: %w", err)
	}

	if animal.Status.IsFinal() {
		return nil, fmt.Errorf("animal.ChangeStatus: already in a final state: %w", domain.ErrValidation)
	}

	newStatus := domain.AnimalStatus(dto.Status)
	if !newStatus.IsFinal() {
		return nil, fmt.Errorf("animal.ChangeStatus: invalid target status: %w", domain.ErrValidation)
	}

	now := time.Now()
	animal.Status = newStatus
	animal.StatusDate = &dto.StatusDate
	
	// Construct structured notes for the status change
	notes := ""
	if dto.StatusNotes != nil {
		notes = *dto.StatusNotes
	}
	if newStatus == domain.StatusSold && dto.BuyerName != nil {
		notes = fmt.Sprintf("Buyer: %s | %s", *dto.BuyerName, notes)
	}
	if newStatus == domain.StatusDead && dto.CauseDeath != nil {
		notes = fmt.Sprintf("Cause: %s | %s", *dto.CauseDeath, notes)
	}
	animal.StatusNotes = &notes
	animal.UpdatedAt = now

	if err := s.repo.Update(ctx, animal); err != nil {
		return nil, fmt.Errorf("animal.ChangeStatus: repository failed: %w", err)
	}

	s.bus.Publish(ports.Event{
		EventType:  "ANIMAL_UPDATED",
		EntityType: "animal",
		EntityID:   animal.ID.String(),
		TenantID:   animal.TenantID.String(),
		Timestamp:  now.Format(time.RFC3339),
	})

	return animal, nil
}
