package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (repo *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
        INSERT INTO subscription (id, tenant_id, sub_source_id, source, subscription_mode, configuration, created_at, active)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = time.Now().UTC()
	}
	if sub.LastPulled.IsZero() {
		sub.LastPulled = time.Now().UTC()
	}

	_, err := repo.db.Exec(ctx, query, sub.ID, sub.TenantID, sub.SubSourceId, sub.Source, sub.SubscriptionMode, sub.Configuration, sub.CreatedAt, true)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %v", err)
	}

	return nil
}

func (repo *SubscriptionRepository) Get(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	query := `SELECT id, tenant_id, sub_source_id, source, subscription_mode, configuration, created_at FROM subscription WHERE id = $1`

	sub := &models.Subscription{}
	row := repo.db.QueryRow(ctx, query, subscriptionID)

	err := row.Scan(&sub.ID, &sub.TenantID, &sub.SubSourceId, &sub.Source, &sub.SubscriptionMode, &sub.Configuration, &sub.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %v", err)
	}

	return sub, nil
}

func (repo *SubscriptionRepository) Delete(ctx context.Context, subscriptionID string) error {
	query := `DELETE FROM subscription WHERE id = $1`

	cmdTag, err := repo.db.Exec(ctx, query, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no subscription found with ID %s", subscriptionID)
	}

	return nil
}

func (repo *SubscriptionRepository) ListByTenantAndApp(ctx context.Context, tenantID, appID string) ([]*models.Subscription, error) {
	query := `SELECT id, tenant_id, sub_source_id, source, subscription_mode, configuration, created_at FROM subscription WHERE tenant_id = $1 AND sub_source_id = $2`

	rows, err := repo.db.Query(ctx, query, tenantID, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscription: %v", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.TenantID, &sub.SubSourceId, &sub.Source, &sub.SubscriptionMode, &sub.Configuration, &sub.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %v", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return subs, nil
}

func (repo *SubscriptionRepository) GetAllActivePullSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	query := `SELECT id, tenant_id, sub_source_id, source, subscription_mode, configuration, created_at FROM subscription WHERE subscription_mode = 'pull' AND active = true`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %v", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.TenantID, &sub.SubSourceId, &sub.Source, &sub.SubscriptionMode, &sub.Configuration, &sub.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %v", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return subs, nil
}

func (repo *SubscriptionRepository) UpdateLastPulled(ctx context.Context, subscriptionID string) error {
	now := time.Now().UTC()
	query := `UPDATE subscription SET last_pulled = $1 WHERE id = $2`
	_, err := repo.db.Exec(ctx, query, now, subscriptionID)
	return err
}
