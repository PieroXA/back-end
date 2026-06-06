package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	healthApp "mifinca-backend/internal/application/health"
	httputil "mifinca-backend/internal/delivery/http/httputil"
	"mifinca-backend/internal/delivery/http/middleware"
)

type HealthHandler struct {
	service *healthApp.Service
	val     *validator.Validate
}

func NewHealthHandler(service *healthApp.Service) *HealthHandler {
	return &HealthHandler{
		service: service,
		val:     validator.New(),
	}
}

type CreateHealthRequest struct {
	EventType   string   `json:"event_type" validate:"required"`
	EventDate   string   `json:"event_date" validate:"required"` // YYYY-MM-DD
	ProductName *string  `json:"product_name,omitempty"`
	DoseMl      *float64 `json:"dose_ml,omitempty"`
	LotNumber   *string  `json:"lot_number,omitempty"`
	VetName     *string  `json:"vet_name,omitempty"`
	NextDate    *string  `json:"next_date,omitempty"` // YYYY-MM-DD
	Notes       *string  `json:"notes,omitempty"`
	Cost        *float64 `json:"cost,omitempty"`
}

func (h *HealthHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	animalIDStr := chi.URLParam(r, "id")
	animalID, err := uuid.Parse(animalIDStr)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid animal id format")
		return
	}

	var req CreateHealthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := h.val.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid event_date format")
		return
	}

	var nextDate *time.Time
	if req.NextDate != nil {
		nd, err := time.Parse("2006-01-02", *req.NextDate)
		if err != nil {
			httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid next_date format")
			return
		}
		nextDate = &nd
	}

	result, err := h.service.Create(r.Context(), healthApp.CreateHealthDTO{
		TenantID:    claims.TenantID,
		AnimalID:    animalID,
		EventType:   req.EventType,
		EventDate:   eventDate,
		ProductName: req.ProductName,
		DoseMl:      req.DoseMl,
		LotNumber:   req.LotNumber,
		VetName:     req.VetName,
		NextDate:    nextDate,
		Notes:       req.Notes,
		Cost:        req.Cost,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, result)
}
