package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type IntercomIntegration struct{}

func NewIntercomStrategy() *IntercomIntegration {
	return &IntercomIntegration{}
}

func (s *IntercomIntegration) Pull(ctx context.Context, sub *models.Subscription) ([]*models.Feedback, error) {
	// pull logic for intercom based on the tenant subscription
	return nil, fmt.Errorf("intercom pull method not implemented yet")
}

func (s *IntercomIntegration) Push(ctx context.Context, r *http.Request, body []byte) ([]*models.Feedback, error) {
	tenantID := r.URL.Query().Get("tenant_id")
	SubSourceID := r.URL.Query().Get("app_id")

	var webhookEvent map[string]interface{}
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook data: %v", err)
	}

	feedback, err := s.processPushRawData(ctx, tenantID, SubSourceID, body)
	if err != nil {
		return nil, fmt.Errorf("failed to process Discourse webhook data: %v", err)
	}

	return []*models.Feedback{feedback}, nil
}

func (a *IntercomIntegration) processPushRawData(ctx context.Context, tenantID string, SubSourceID string, data []byte) (*models.Feedback, error) {
	var intercomData struct {
		ID             string `json:"id"`
		ConversationID string `json:"conversation_id"`
		Messages       []struct {
			ID   string `json:"id"`
			Body string `json:"body"`
		} `json:"conversation_parts"`
	}

	if err := json.Unmarshal(data, &intercomData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Intercom data: %v", err)
	}

	messages := make([]models.Message, len(intercomData.Messages))
	for i, msg := range intercomData.Messages {
		messages[i] = models.Message{
			ID:      msg.ID,
			Content: msg.Body,
		}
	}

	content := models.ConversationContent{
		ConversationID: intercomData.ConversationID,
		Messages:       messages,
	}

	// Create and return the FeedbackRecord
	return &models.Feedback{
		ID:          intercomData.ID,
		TenantID:    tenantID,
		SubSourceID: SubSourceID,
		Source:      a.GetSourceName(),
		Type:        "Conversation",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    map[string]interface{}{},
		Content:     content,
	}, nil
}

func (a *IntercomIntegration) GetSourceName() models.Source {
	return models.SourceIntercom
}

func (a *IntercomIntegration) GetSourceType() models.SourceType {
	return models.STConversation
}
