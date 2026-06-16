package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

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

	// For now, let's assume Host is used if Origin is missing.
	rawDomain := r.Header.Get("Origin")
	if rawDomain == "" {
		rawDomain = r.Host
	}

	domain := normalizeDomain(rawDomain)

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

func normalizeDomain(raw string) string {
	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err == nil {
			host := u.Hostname()
			if host != "" {
				return host
			}
		}
	}
	host := raw
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
	}
	return host
}
