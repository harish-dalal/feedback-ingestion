package server

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateTables(ctx context.Context, dbpool *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tenant (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			api_key TEXT NOT NULL UNIQUE
		);`,

		`CREATE TABLE IF NOT EXISTS feedback (
			id TEXT NOT NULL,
			tenant_id UUID NOT NULL,
			sub_source_id UUID NOT NULL,
			source TEXT NOT NULL,
			source_type TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			metadata JSONB,
			content JSONB,
			CONSTRAINT pk_feedback PRIMARY KEY (id, tenant_id, source),
			CONSTRAINT fk_tenant
			  FOREIGN KEY(tenant_id) 
			  REFERENCES tenant(id)
			  ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS subscription (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL,
			sub_source_id UUID NOT NULL,
			source TEXT NOT NULL,
			subscription_mode TEXT NOT NULL CHECK (subscription_mode IN ('push', 'pull')),
			configuration JSONB,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			last_pulled TIMESTAMPTZ,
			active BOOLEAN NOT NULL DEFAULT TRUE,
			CONSTRAINT fk_tenant
			  FOREIGN KEY(tenant_id) 
			  REFERENCES tenant(id)
			  ON DELETE CASCADE
		);`,

		// Insert a default tenant if not already exists
		`INSERT INTO tenant (id, name, api_key)
		 VALUES ('cb4d81c7-e1bf-4ca5-900f-665a0e3fc932', 'Default Tenant', '3eaa4c4a-2271-45b1-8df0-81a4bd64251b')
		 ON CONFLICT (id) DO NOTHING;`,
	}

	for _, query := range queries {
		if _, err := dbpool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
