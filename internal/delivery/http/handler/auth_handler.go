package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	authApp "mifinca-backend/internal/application/auth"
	httputil "mifinca-backend/internal/delivery/http/httputil"
)

var validate = validator.New()

type AuthHandler struct {
	service *authApp.Service
}

func NewAuthHandler(service *authApp.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// ── Request types ─────────────────────────────────────────────────────────────

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=200"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
	FarmName string `json:"farm_name" validate:"required,min=2,max=200"`
}

type LoginRequest struct {
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required"`
	DeviceID string `json:"device_id" validate:"required,uuid4"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	result, err := h.service.Register(r.Context(), authApp.RegisterDTO{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		FarmName: req.FarmName,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	result, err := h.service.Login(r.Context(), authApp.LoginDTO{
		Email:    req.Email,
		Password: req.Password,
		DeviceID: req.DeviceID,
	})
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	result, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, result)
}
