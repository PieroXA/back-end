package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	animalApp "mifinca-backend/internal/application/animal"
	httputil "mifinca-backend/internal/delivery/http/httputil"
	"mifinca-backend/internal/delivery/http/middleware"
	"mifinca-backend/internal/domain/ports"
)

type AnimalHandler struct {
	service *animalApp.Service
}

func NewAnimalHandler(service *animalApp.Service) *AnimalHandler {
	return &AnimalHandler{service: service}
}

// ── Request Types ─────────────────────────────────────────────────────────────

type CreateAnimalRequest struct {
	VisualTag string   `json:"visual_tag" validate:"required,max=50"`
	RfidTag   *string  `json:"rfid_tag,omitempty" validate:"omitempty,max=100"`
	Species   string   `json:"species"    validate:"required,oneof=bovino ovino caprino porcino equino"`
	Breed     *string  `json:"breed,omitempty"`
	Sex       string   `json:"sex"        validate:"required,oneof=M F"`
	BirthDate string   `json:"birth_date" validate:"required"` // YYYY-MM-DD
	WeightKg  *float64 `json:"weight_kg,omitempty"  validate:"omitempty,gt=0"`
	Color     *string  `json:"color,omitempty"`
	SireID    *string  `json:"sire_id,omitempty"    validate:"omitempty,uuid4"`
	DamID     *string  `json:"dam_id,omitempty"     validate:"omitempty,uuid4"`
	Notes     *string  `json:"notes,omitempty"`
}

type UpdateAnimalRequest struct {
	VisualTag *string  `json:"visual_tag,omitempty" validate:"omitempty,max=50"`
	RfidTag   *string  `json:"rfid_tag,omitempty" validate:"omitempty,max=100"`
	Breed     *string  `json:"breed,omitempty"`
	WeightKg  *float64 `json:"weight_kg,omitempty" validate:"omitempty,gt=0"`
	Color     *string  `json:"color,omitempty"`
	SireID    *string  `json:"sire_id,omitempty" validate:"omitempty,uuid4"`
	DamID     *string  `json:"dam_id,omitempty" validate:"omitempty,uuid4"`
	FarmID    *string  `json:"farm_id,omitempty" validate:"omitempty,uuid4"`
	Notes     *string  `json:"notes,omitempty"`
}

type ChangeStatusRequest struct {
	Status      string  `json:"status"      validate:"required,oneof=sold dead lost"`
	StatusDate  string  `json:"status_date" validate:"required"` // YYYY-MM-DD
	StatusNotes *string `json:"status_notes,omitempty"`
	BuyerName   *string `json:"buyer_name,omitempty"`
	CauseDeath  *string `json:"cause_of_death,omitempty"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func (h *AnimalHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	farmIDStr := chi.URLParam(r, "farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid farm_id format")
		return
	}

	var req CreateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid birth_date format (use YYYY-MM-DD)")
		return
	}

	var sireID, damID *uuid.UUID
	if req.SireID != nil {
		id, _ := uuid.Parse(*req.SireID)
		sireID = &id
	}
	if req.DamID != nil {
		id, _ := uuid.Parse(*req.DamID)
		damID = &id
	}

	result, err := h.service.Create(r.Context(), animalApp.CreateAnimalDTO{
		TenantID:  claims.TenantID,
		FarmID:    farmID,
		VisualTag: req.VisualTag,
		RfidTag:   req.RfidTag,
		Species:   req.Species,
		Breed:     req.Breed,
		Sex:       req.Sex,
		BirthDate: birthDate,
		WeightKg:  req.WeightKg,
		Color:     req.Color,
		SireID:    sireID,
		DamID:     damID,
		Notes:     req.Notes,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, result)
}

func (h *AnimalHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	var req UpdateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	dto := animalApp.UpdateAnimalDTO{
		TenantID:  claims.TenantID,
		AnimalID:  animalID,
		VisualTag: req.VisualTag,
		RfidTag:   req.RfidTag,
		Breed:     req.Breed,
		WeightKg:  req.WeightKg,
		Color:     req.Color,
		Notes:     req.Notes,
	}

	if req.SireID != nil {
		id, _ := uuid.Parse(*req.SireID)
		dto.SireID = &id
	}
	if req.DamID != nil {
		id, _ := uuid.Parse(*req.DamID)
		dto.DamID = &id
	}
	if req.FarmID != nil {
		id, _ := uuid.Parse(*req.FarmID)
		dto.FarmID = &id
	}

	result, err := h.service.Update(r.Context(), dto)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}

func (h *AnimalHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	var req ChangeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	statusDate, err := time.Parse("2006-01-02", req.StatusDate)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid status_date format")
		return
	}

	result, err := h.service.ChangeStatus(r.Context(), animalApp.ChangeStatusDTO{
		TenantID:    claims.TenantID,
		AnimalID:    animalID,
		Status:      req.Status,
		StatusDate:  statusDate,
		StatusNotes: req.StatusNotes,
		BuyerName:   req.BuyerName,
		CauseDeath:  req.CauseDeath,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}

func (h *AnimalHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	if err := h.service.SoftDelete(r.Context(), claims.TenantID, animalID); err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusNoContent, nil)
}

func (h *AnimalHandler) GetPedigree(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	result, err := h.service.GetPedigree(r.Context(), claims.TenantID, animalID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}

func (h *AnimalHandler) GetAnimal(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	result, err := h.service.GetByID(r.Context(), claims.TenantID, animalID)
	if err != nil {
		httputil.Error(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}

func (h *AnimalHandler) ListAnimals(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	
	farmIDStr := chi.URLParam(r, "farm_id")
	var farmID *uuid.UUID
	if farmIDStr != "" {
		id, err := uuid.Parse(farmIDStr)
		if err == nil {
			farmID = &id
		}
	}

	q := r.URL.Query()
	
	var species, sex, status, search, cursor *string
	if v := q.Get("species"); v != "" { species = &v }
	if v := q.Get("sex"); v != "" { sex = &v }
	if v := q.Get("status"); v != "" { status = &v }
	if v := q.Get("search"); v != "" { search = &v }
	if v := q.Get("cursor"); v != "" { cursor = &v }

	limitStr := q.Get("limit")
	limit := int32(50)
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = int32(l)
		}
	}

	params := ports.ListAnimalsParams{
		TenantID: claims.TenantID,
		FarmID:   farmID,
		Species:  species,
		Sex:      sex,
		Status:   status,
		Search:   search,
		Limit:    limit,
		Cursor:   cursor,
	}

	animals, count, err := h.service.List(r.Context(), params)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	var nextCursor *string
	if len(animals) > 0 && int32(len(animals)) == limit {
		lastTag := animals[len(animals)-1].VisualTag
		nextCursor = &lastTag
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"data": animals,
		"total": count,
		"next_cursor": nextCursor,
	})
}
