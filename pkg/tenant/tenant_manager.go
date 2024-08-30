package tenant

import (
	"errors"
	"sync"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type TenantManager struct {
	tenants map[string]*models.Tenant
	mutex   sync.RWMutex
}

func NewTenantManager() *TenantManager {
	return &TenantManager{
		tenants: make(map[string]*models.Tenant),
	}
}

func (tm *TenantManager) AddTenant(tenant *models.Tenant) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.tenants[tenant.ID] = tenant
}

func (tm *TenantManager) GetTenant(tenantID string) (*models.Tenant, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	tenant, exists := tm.tenants[tenantID]
	if !exists {
		return nil, errors.New("tenant not found")
	}
	return tenant, nil
}

func (tm *TenantManager) GetTenantByApiKey(apiKey string) (*models.Tenant, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	for _, tenant := range tm.tenants {
		if tenant.ApiKey == apiKey {
			return tenant, nil
		}
	}
	return nil, errors.New("tenant not found")
}
