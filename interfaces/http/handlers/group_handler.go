package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/user"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/interfaces/http/middleware"
	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	service user.GroupService
}

func NewGroupHandler(service user.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req user.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.TenantID = tenantID

	res, err := h.service.CreateGroup(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	res, err := h.service.ListGroups(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *GroupHandler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")

	if err := h.service.AddUserToGroup(r.Context(), groupID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *GroupHandler) AddRoleToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "roleId")

	if err := h.service.AddRoleToGroup(r.Context(), groupID, roleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
