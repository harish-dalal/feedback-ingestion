package db

import (
	"context"
	"fmt"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TenantRepository struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{db: db}
}

func (repo *TenantRepository) Save(ctx context.Context, tenant *models.Tenant) error {
	query := `
        INSERT INTO tenants (id, name, api_key)
        VALUES ($1, $2, $3)
    `

	_, err := repo.db.Exec(ctx, query, tenant.ID, tenant.Name, tenant.ApiKey)
	if err != nil {
		return fmt.Errorf("failed to save tenant: %v", err)
	}

	return nil
}

func (repo *TenantRepository) Get(ctx context.Context, tenantID string) (*models.Tenant, error) {
	query := `SELECT id, name, api_key FROM tenants WHERE id = $1`

	tenant := &models.Tenant{}
	row := repo.db.QueryRow(ctx, query, tenantID)

	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}

	return tenant, nil
}

func (repo *TenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	query := `
        UPDATE tenants
        SET name = $2
        WHERE id = $1
    `

	_, err := repo.db.Exec(ctx, query, tenant.ID, tenant.Name)
	if err != nil {
		return fmt.Errorf("failed to update tenant: %v", err)
	}

	return nil
}

func (repo *TenantRepository) Delete(ctx context.Context, tenantID string) error {
	query := `DELETE FROM tenants WHERE id = $1`

	_, err := repo.db.Exec(ctx, query, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %v", err)
	}

	return nil
}
