package feedback

import (
	"context"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/db"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type FeedbackService struct {
	repo *db.FeedbackRepository
}

func NewFeedbackService(repo *db.FeedbackRepository) *FeedbackService {
	return &FeedbackService{repo: repo}
}

func (s *FeedbackService) CreateFeedback(ctx context.Context, feedback *models.Feedback) error {
	return s.repo.Save(ctx, feedback)
}

func (s *FeedbackService) GetFeedback(ctx context.Context, feedbackID string) (*models.Feedback, error) {
	return s.repo.Get(ctx, feedbackID)
}

func (s *FeedbackService) UpdateFeedback(ctx context.Context, feedback *models.Feedback) error {
	return s.repo.Update(ctx, feedback)
}

func (s *FeedbackService) DeleteFeedback(ctx context.Context, feedbackID string) error {
	return s.repo.Delete(ctx, feedbackID)
}

func (s *FeedbackService) ListFeedbackByTenant(ctx context.Context, tenantID string) ([]*models.Feedback, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}
