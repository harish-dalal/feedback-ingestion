package tenant

import (
	"context"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/db"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// TenantService provides methods for managing tenants
type TenantService struct {
	repo *db.TenantRepository
}

// NewTenantService creates a new TenantService
func NewTenantService(repo *db.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

// CreateTenant creates a new tenant
func (s *TenantService) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.Save(ctx, tenant)
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	return s.repo.Get(ctx, tenantID)
}

// UpdateTenant updates an existing tenant
func (s *TenantService) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.Update(ctx, tenant)
}

// DeleteTenant deletes a tenant by ID
func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	return s.repo.Delete(ctx, tenantID)
}
