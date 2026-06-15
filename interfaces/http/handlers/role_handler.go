package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/role"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	service role.Service
}

func NewRoleHandler(service role.Service) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req role.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.TenantID = tenantID

	res, err := h.service.CreateRole(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.service.GetRole(r.Context(), id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	res, err := h.service.ListRoles(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req role.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.ID = id

	res, err := h.service.UpdateRole(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteRole(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
