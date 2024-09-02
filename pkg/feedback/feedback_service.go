package feedback

import (
	"context"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/db"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// FeedbackService provides methods for managing feedback records
type FeedbackService struct {
	repo *db.FeedbackRepository
}

// NewFeedbackService creates a new FeedbackService
func NewFeedbackService(repo *db.FeedbackRepository) *FeedbackService {
	return &FeedbackService{repo: repo}
}

// CreateFeedback creates a new feedback record
func (s *FeedbackService) CreateFeedback(ctx context.Context, record *models.Feedback) error {
	return s.repo.Save(ctx, record)
}

// GetFeedback retrieves a feedback record by ID
func (s *FeedbackService) GetFeedback(ctx context.Context, feedbackID string) (*models.Feedback, error) {
	return s.repo.Get(ctx, feedbackID)
}

// UpdateFeedback updates an existing feedback record
func (s *FeedbackService) UpdateFeedback(ctx context.Context, record *models.Feedback) error {
	return s.repo.Update(ctx, record)
}

// DeleteFeedback deletes a feedback record by ID
func (s *FeedbackService) DeleteFeedback(ctx context.Context, feedbackID string) error {
	return s.repo.Delete(ctx, feedbackID)
}

// ListFeedbackByTenant lists all feedback records for a given tenant
func (s *FeedbackService) ListFeedbackByTenant(ctx context.Context, tenantID string) ([]*models.Feedback, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}
