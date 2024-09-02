package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SubscriptionRepository defines methods for managing subscriptions in the PostgreSQL database
type SubscriptionRepository struct {
	db *pgxpool.Pool
}

// NewSubscriptionRepository creates a new SubscriptionRepository
func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription in the database
func (repo *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
        INSERT INTO subscriptions (id, tenant_id, app_id, source_id, subscription_mode, configuration, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = time.Now().UTC()
	}

	_, err := repo.db.Exec(ctx, query, sub.ID, sub.TenantID, sub.AppID, sub.SourceID, sub.SubscriptionMode, sub.Configuration, sub.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %v", err)
	}

	return nil
}

// Get retrieves a subscription by ID from the database
func (repo *SubscriptionRepository) Get(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	query := `SELECT id, tenant_id, app_id, source_id, subscription_mode, configuration, created_at FROM subscriptions WHERE id = $1`

	sub := &models.Subscription{}
	row := repo.db.QueryRow(ctx, query, subscriptionID)

	err := row.Scan(&sub.ID, &sub.TenantID, &sub.AppID, &sub.SourceID, &sub.SubscriptionMode, &sub.Configuration, &sub.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %v", err)
	}

	return sub, nil
}

// Delete deletes a subscription by ID from the database
func (repo *SubscriptionRepository) Delete(ctx context.Context, subscriptionID string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	cmdTag, err := repo.db.Exec(ctx, query, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no subscription found with ID %s", subscriptionID)
	}

	return nil
}

// ListByTenantAndApp retrieves all subscriptions for a given tenant and app
func (repo *SubscriptionRepository) ListByTenantAndApp(ctx context.Context, tenantID, appID string) ([]*models.Subscription, error) {
	query := `SELECT id, tenant_id, app_id, source_id, subscription_mode, configuration, created_at FROM subscriptions WHERE tenant_id = $1 AND app_id = $2`

	rows, err := repo.db.Query(ctx, query, tenantID, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %v", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.TenantID, &sub.AppID, &sub.SourceID, &sub.SubscriptionMode, &sub.Configuration, &sub.CreatedAt); err != nil {
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
	// Get all active pull subscriptions from the database
	query := `SELECT id, tenant_id, app_id, source_id, subscription_mode, configuration, created_at FROM subscriptions WHERE subscription_mode = 'pull'`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %v", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.TenantID, &sub.AppID, &sub.SourceID, &sub.SubscriptionMode, &sub.Configuration, &sub.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %v", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return subs, nil
}

// UpdateLastPulled updates the last_pulled timestamp for a subscription.
func (repo *SubscriptionRepository) UpdateLastPulled(ctx context.Context, subscriptionID string) error {
	now := time.Now().UTC()
	query := `UPDATE subscriptions SET last_pulled = $1 WHERE id = $2`
	_, err := repo.db.Exec(ctx, query, now, subscriptionID)
	return err
}
