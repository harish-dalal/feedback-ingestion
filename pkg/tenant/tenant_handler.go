package tenant

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// TenantHandler provides HTTP handlers for managing tenants
type TenantHandler struct {
	service *TenantService
}

// NewTenantHandler creates a new TenantHandler
func NewTenantHandler(service *TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

// CreateTenantHandler handles HTTP requests to create a new tenant
func (h *TenantHandler) CreateTenantHandler(w http.ResponseWriter, r *http.Request) {
	var tenant models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tenant.ID = uuid.New().String()
	tenant.ApiKey = uuid.New().String()

	ctx := r.Context()
	if err := h.service.CreateTenant(ctx, &tenant); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create tenant: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}

// GetTenantHandler handles HTTP requests to retrieve a tenant by ID
func (h *TenantHandler) GetTenantHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tenant, err := h.service.GetTenant(ctx, tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve tenant: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}

// UpdateTenantHandler handles HTTP requests to update an existing tenant
func (h *TenantHandler) UpdateTenantHandler(w http.ResponseWriter, r *http.Request) {
	var tenant models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if tenant.ID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.UpdateTenant(ctx, &tenant); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update tenant: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}

// DeleteTenantHandler handles HTTP requests to delete a tenant by ID
func (h *TenantHandler) DeleteTenantHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteTenant(ctx, tenantID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete tenant: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
