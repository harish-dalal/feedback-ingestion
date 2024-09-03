package server

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateTables(ctx context.Context, dbpool *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			api_key TEXT NOT NULL UNIQUE
		);`,

		`CREATE TABLE IF NOT EXISTS feedback (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL,
			sub_source_id UUID NOT NULL,
			source TEXT NOT NULL,
			type TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			metadata JSONB,
			content JSONB,
			CONSTRAINT fk_tenant
			  FOREIGN KEY(tenant_id) 
			  REFERENCES tenants(id)
			  ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS subscriptions (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL,
			app_id UUID NOT NULL,
			source TEXT NOT NULL,
			subscription_mode TEXT NOT NULL CHECK (subscription_mode IN ('push', 'pull')),
			configuration JSONB,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			last_pulled TIMESTAMPTZ,
			active BOOLEAN NOT NULL DEFAULT TRUE,
			CONSTRAINT fk_tenant
			  FOREIGN KEY(tenant_id) 
			  REFERENCES tenants(id)
			  ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		if _, err := dbpool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
