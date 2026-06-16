package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/tenant"
)

type TenantHandler struct {
	service tenant.Service
}

func NewTenantHandler(service tenant.Service) *TenantHandler {
	return &TenantHandler{service: service}
}

func (h *TenantHandler) RegisterTenant(w http.ResponseWriter, r *http.Request) {
	var req tenant.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.RegisterTenant(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.service.ListTenants(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}
