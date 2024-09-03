package integrations

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/feedback"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type SourceStrategy interface {
	// Pull - pulls the data from the source and saves it to the feedback database
	Pull(ctx context.Context, sub *models.Subscription) ([]*models.Feedback, error)

	// Push - recives the data from source's webhook and saves it to the feedback database
	Push(ctx context.Context, r *http.Request, body []byte) ([]*models.Feedback, error)

	// GetSourceName ...
	GetSourceName() models.Source

	// GetSourceType ...
	GetSourceType() models.SourceType
}

type IntegrationManager struct {
	strategies      map[models.Source]SourceStrategy
	feedbackService *feedback.FeedbackService
}

func NewIntegrationManager(strategies map[models.Source]SourceStrategy, feedbackService *feedback.FeedbackService) *IntegrationManager {
	return &IntegrationManager{strategies: strategies, feedbackService: feedbackService}
}

func (m *IntegrationManager) Pull(ctx context.Context, sub *models.Subscription) ([]*models.Feedback, error) {
	strategy, ok := m.strategies[models.Source(sub.Source)]
	if !ok {
		return nil, fmt.Errorf("no strategy found for source: %s", sub.Source)
	}
	feedbacks, err := strategy.Pull(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("failed to pull data from source: %v", err)
	}

	err = m.feedbacksToDB(feedbacks)
	if err != nil {
		// log error
	}

	return feedbacks, nil
}

func (m *IntegrationManager) HandleWebhook(w http.ResponseWriter, r *http.Request, source models.Source) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	strategy, ok := m.strategies[source]
	if !ok {
		http.Error(w, fmt.Sprintf("no strategy found for source: %s", strategy), http.StatusInternalServerError)
		return
	}

	feedbacks, err := strategy.Push(ctx, r, body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process webhook: %v", err), http.StatusInternalServerError)
		return
	}

	err = m.feedbacksToDB(feedbacks)
	if err != nil {
		// log the error
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Discourse webhook received successfully")
}

func (m *IntegrationManager) feedbacksToDB(feedbacks []*models.Feedback) error {
	for _, feedback := range feedbacks {
		err := m.feedbackService.CreateFeedback(context.Background(), feedback)
		if err != nil {
			// log error;
		}
	}
	return nil
}
