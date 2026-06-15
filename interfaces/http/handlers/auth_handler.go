package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/application/auth"
)

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Strategy    string            `json:"strategy"`
	Credentials map[string]string `json:"credentials"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Domain extraction for Feature A is handled by middleware, 
	// but we can also extract it here if needed or use context.
	// For now, let's assume Host is used if Origin is missing.
	domain := r.Header.Get("Origin")
	if domain == "" {
		domain = r.Host
	}

	token, err := h.authService.Authenticate(r.Context(), domain, req.Strategy, req.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}
