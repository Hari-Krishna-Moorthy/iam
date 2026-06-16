package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

type BulkHandler struct {
	bulkService user.BulkService
}

func NewBulkHandler(bulkService user.BulkService) *BulkHandler {
	return &BulkHandler{bulkService: bulkService}
}

func (h *BulkHandler) BulkCreateUsers(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Tenant not found", http.StatusUnauthorized)
		return
	}

	var req user.BulkCreateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	jobID, err := h.bulkService.SubmitBulkCreate(r.Context(), tenantID, req)
	if err != nil {
		http.Error(w, "Failed to submit job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted for async jobs
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID, "status": "pending"})
}

func (h *BulkHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	j, err := h.bulkService.GetJobStatus(r.Context(), jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// In a real app, ensure the user requesting the job status belongs to the tenant that owns the job.
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	if j.TenantID != tenantID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}
