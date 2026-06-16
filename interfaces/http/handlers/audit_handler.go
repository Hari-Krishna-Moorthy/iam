package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/audit"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
)

type AuditHandler struct {
	repo audit.Repository
}

func NewAuditHandler(repo audit.Repository) *AuditHandler {
	return &AuditHandler{repo: repo}
}

func (h *AuditHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Tenant not found", http.StatusUnauthorized)
		return
	}

	logs, err := h.repo.GetByTenantID(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
