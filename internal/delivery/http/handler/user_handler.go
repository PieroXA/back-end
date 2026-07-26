package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"mifinca-backend/internal/application/user"
	"mifinca-backend/internal/delivery/http/httputil"
	"mifinca-backend/internal/delivery/http/middleware"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
		return
	}

	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}

	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	res, err := h.userService.Create(r.Context(), claims.TenantID, req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, res)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
		return
	}

	res, err := h.userService.List(r.Context(), claims.TenantID)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, res)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid user ID")
		return
	}

	res, err := h.userService.GetByID(r.Context(), claims.TenantID, id)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, res)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid user ID")
		return
	}

	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid JSON body")
		return
	}

	if err := validate.Struct(req); err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}

	res, err := h.userService.Update(r.Context(), claims.TenantID, id, req)
	if err != nil {
		httputil.Error(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, res)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httputil.ErrorWithMessage(w, http.StatusBadRequest, "VALIDATION", "invalid user ID")
		return
	}

	if err := h.userService.Delete(r.Context(), claims.TenantID, id); err != nil {
		httputil.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
