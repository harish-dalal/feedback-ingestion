package integrations

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type SourceStrategy interface {
	Pull(ctx context.Context, sub *models.Subscription) error
	Push(ctx context.Context, r *http.Request, body []byte) error
	ProcessRawData(ctx context.Context, tenantID string, data []byte) (*models.Feedback, error)
}

type IntegrationManager struct {
	strategies map[models.SourceType]SourceStrategy
}

func NewIntegrationManager(strategies map[models.SourceType]SourceStrategy) *IntegrationManager {
	return &IntegrationManager{strategies: strategies}
}

func (m *IntegrationManager) Pull(ctx context.Context, sub *models.Subscription) error {
	strategy, ok := m.strategies[models.SourceType(sub.SourceID)]
	if !ok {
		return fmt.Errorf("no strategy found for source: %s", sub.SourceID)
	}
	return strategy.Pull(ctx, sub)
}

func (m *IntegrationManager) HandleWebhook(w http.ResponseWriter, r *http.Request, source models.SourceType) {
	ctx := r.Context()

	// Step 1: Parse the incoming data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Step 2: Select strategy and push
	strategy, ok := m.strategies[source]
	if !ok {
		http.Error(w, fmt.Sprintf("no strategy found for source: %s", strategy), http.StatusInternalServerError)
		return
	}

	err = strategy.Push(ctx, r, body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process webhook: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 3: Respond to the webhook request
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Discourse webhook received successfully")
}

// func (m *IntegrationManager) Push(ctx context.Context, r *http.Request, body []byte) error {
// 	strategy, ok := m.strategies[source]
// 	if !ok {
// 		return fmt.Errorf("no strategy found for source: %s", strategy)
// 	}
// 	return strategy.Push(ctx, r, body, source)
// }

func (m *IntegrationManager) ProcessData(ctx context.Context, source models.SourceType, tenantID string, data []byte) (*models.Feedback, error) {
	strategy, ok := m.strategies[source]
	if !ok {
		return nil, fmt.Errorf("no strategy found for source: %s", source)
	}
	return strategy.ProcessRawData(ctx, tenantID, data)
}
