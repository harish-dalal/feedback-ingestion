package subscription

import (
	"context"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/db"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type SubscriptionService struct {
	repo *db.SubscriptionRepository
}

func NewSubscriptionService(repo *db.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	if err := s.repo.Create(ctx, sub); err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) GetAllActivePullSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	subscriptions, err := s.repo.GetAllActivePullSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionService) UpdateLastPulled(ctx context.Context, subscriptionID string) error {
	err := s.repo.UpdateLastPulled(ctx, subscriptionID)
	if err != nil {
		return err
	}
	return nil
}
