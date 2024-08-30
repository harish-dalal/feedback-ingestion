package db

import (
	"context"
	"fmt"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TenantRepository defines methods for managing tenants in the PostgreSQL database
type TenantRepository struct {
	db *pgxpool.Pool
}

// NewTenantRepository creates a new TenantRepository
func NewTenantRepository(db *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{db: db}
}

// Save saves a new tenant to the database
func (repo *TenantRepository) Save(ctx context.Context, tenant *models.Tenant) error {
	query := `
        INSERT INTO tenants (id, name, api_key, active_sources, configurations)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := repo.db.Exec(ctx, query, tenant.ID, tenant.Name, tenant.ApiKey, tenant.ActiveSources, tenant.Configurations)
	if err != nil {
		return fmt.Errorf("failed to save tenant: %v", err)
	}

	return nil
}

// Get retrieves a tenant by ID from the database
func (repo *TenantRepository) Get(ctx context.Context, tenantID string) (*models.Tenant, error) {
	query := `SELECT id, name, api_key, active_sources, configurations FROM tenants WHERE id = $1`

	tenant := &models.Tenant{}
	row := repo.db.QueryRow(ctx, query, tenantID)

	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.ApiKey, &tenant.ActiveSources, &tenant.Configurations)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}

	return tenant, nil
}

// Update updates an existing tenant in the database
func (repo *TenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	query := `
        UPDATE tenants
        SET name = $2, active_sources = $3, configurations = $4
        WHERE id = $1
    `

	_, err := repo.db.Exec(ctx, query, tenant.ID, tenant.Name, tenant.ActiveSources, tenant.Configurations)
	if err != nil {
		return fmt.Errorf("failed to update tenant: %v", err)
	}

	return nil
}

// Delete deletes a tenant by ID from the database
func (repo *TenantRepository) Delete(ctx context.Context, tenantID string) error {
	query := `DELETE FROM tenants WHERE id = $1`

	_, err := repo.db.Exec(ctx, query, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %v", err)
	}

	return nil
}
