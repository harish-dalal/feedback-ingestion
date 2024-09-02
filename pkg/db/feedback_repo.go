package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FeedbackRepository struct {
	db *pgxpool.Pool
}

func NewFeedbackRepository(db *pgxpool.Pool) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (repo *FeedbackRepository) Save(ctx context.Context, feedback *models.Feedback) error {
	query := `
        INSERT INTO feedback_records (id, tenant_id, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	if feedback.ID == "" {
		feedback.ID = uuid.New().String()
	}
	if feedback.CreatedAt.IsZero() {
		feedback.CreatedAt = time.Now().UTC()
	}
	feedback.UpdatedAt = time.Now().UTC()

	_, err := repo.db.Exec(ctx, query, feedback.ID, feedback.TenantID, feedback.Content, feedback.CreatedAt, feedback.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save feedback %v", err)
	}

	return nil
}

// Get retrieves a feedback record by ID from the database
func (repo *FeedbackRepository) Get(ctx context.Context, feedbackID string) (*models.Feedback, error) {
	query := `SELECT id, tenant_id, content, created_at, updated_at FROM feedback_records WHERE id = $1`

	record := &models.Feedback{}
	row := repo.db.QueryRow(ctx, query, feedbackID)

	err := row.Scan(&record.ID, &record.TenantID, &record.Content, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback record: %v", err)
	}

	return record, nil
}

// Update updates an existing feedback record in the database
func (repo *FeedbackRepository) Update(ctx context.Context, feedback *models.Feedback) error {
	query := `
        UPDATE feedback_records
        SET content = $2, updated_at = $3
        WHERE id = $1
    `
	feedback.UpdatedAt = time.Now().UTC()

	cmdTag, err := repo.db.Exec(ctx, query, feedback.ID, feedback.Content, feedback.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update feedback record: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no feedback record found with ID %s", feedback.ID)
	}

	return nil
}

// Delete deletes a feedback record by ID from the database
func (repo *FeedbackRepository) Delete(ctx context.Context, feedbackID string) error {
	query := `DELETE FROM feedback_records WHERE id = $1`

	cmdTag, err := repo.db.Exec(ctx, query, feedbackID)
	if err != nil {
		return fmt.Errorf("failed to delete feedback record: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no feedback record found with ID %s", feedbackID)
	}

	return nil
}

// ListByTenant retrieves all feedback records for a given tenant
func (repo *FeedbackRepository) ListByTenant(ctx context.Context, tenantID string) ([]*models.Feedback, error) {
	query := `SELECT id, tenant_id, content, created_at, updated_at FROM feedback_records WHERE tenant_id = $1 ORDER BY created_at DESC`

	rows, err := repo.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedback records: %v", err)
	}
	defer rows.Close()

	var records []*models.Feedback
	for rows.Next() {
		record := &models.Feedback{}
		if err := rows.Scan(&record.ID, &record.TenantID, &record.Content, &record.CreatedAt, &record.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan feedback record: %v", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return records, nil
}
