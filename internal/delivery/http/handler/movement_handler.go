package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	movementApp "mifinca-backend/internal/application/movement"
	httputil "mifinca-backend/internal/delivery/http/httputil"
	"mifinca-backend/internal/delivery/http/middleware"
)

type MovementHandler struct {
	service *movementApp.Service
	val     *validator.Validate
}

func NewMovementHandler(service *movementApp.Service) *MovementHandler {
	return &MovementHandler{
		service: service,
		val:     validator.New(),
	}
}

type CreateMovementRequest struct {
	FromFarmID   *string `json:"from_farm_id,omitempty" validate:"omitempty,uuid4"`
	ToFarmID     *string `json:"to_farm_id,omitempty"   validate:"omitempty,uuid4"`
	MovementDate string  `json:"movement_date"          validate:"required"` // YYYY-MM-DD
	Reason       string  `json:"reason"                 validate:"required"`
	Notes        *string `json:"notes,omitempty"`
}

func (h *MovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	var req CreateMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := h.val.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	mDate, err := time.Parse("2006-01-02", req.MovementDate)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid movement_date format")
		return
	}

	var fromFarm, toFarm *uuid.UUID
	if req.FromFarmID != nil {
		id, _ := uuid.Parse(*req.FromFarmID)
		fromFarm = &id
	}
	if req.ToFarmID != nil {
		id, _ := uuid.Parse(*req.ToFarmID)
		toFarm = &id
	}

	result, err := h.service.Create(r.Context(), movementApp.CreateMovementDTO{
		TenantID:     claims.TenantID,
		AnimalID:     animalID,
		FromFarmID:   fromFarm,
		ToFarmID:     toFarm,
		MovementDate: mDate,
		Reason:       req.Reason,
		Notes:        req.Notes,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, result)
}
