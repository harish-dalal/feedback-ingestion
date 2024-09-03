package tenant

import (
	"context"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/db"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type TenantService struct {
	repo *db.TenantRepository
}

func NewTenantService(repo *db.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.Save(ctx, tenant)
}

func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	return s.repo.Get(ctx, tenantID)
}

func (s *TenantService) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.repo.Update(ctx, tenant)
}

func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	return s.repo.Delete(ctx, tenantID)
}
