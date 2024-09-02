package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// IntercomAdapter implements the SourceAdapter interface for Intercom data
type IntercomIntegration struct{}

// NewIntercomStrategy creates a new instance of IntercomStrategy.
func NewIntercomStrategy() *IntercomIntegration {
	return &IntercomIntegration{}
}

func (s *IntercomIntegration) Pull(ctx context.Context, sub *models.Subscription) error {
	lastPulled := sub.LastPulled.Format(time.RFC3339)

	url := fmt.Sprintf("https://api.intercom.io/conversations?since=%s", lastPulled)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+sub.Configuration["access_token"].(string))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	var data []byte
	_, err = resp.Body.Read(data)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Process the pulled data
	// You would typically call ProcessRawData here to convert the raw data to a FeedbackRecord
	fmt.Println("Successfully pulled data from Discourse")
	return nil
}

func (s *IntercomIntegration) Push(ctx context.Context, r *http.Request, body []byte) error {
	// Extract tenant ID and app ID from query parameters
	tenantID := r.URL.Query().Get("tenant_id")
	// appID := r.URL.Query().Get("app_id")

	// Transform the webhook data into your internal FeedbackRecord model
	var webhookEvent map[string]interface{}
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		return fmt.Errorf("failed to unmarshal webhook data: %v", err)
	}

	// Call the IntegrationManager to process the data for Discourse
	_, err := s.ProcessRawData(ctx, tenantID, body)
	if err != nil {
		return fmt.Errorf("failed to process Discourse webhook data: %v", err)
	}

	return nil
}

// ProcessRawData processes raw Intercom data and returns a FeedbackRecord
func (a *IntercomIntegration) ProcessRawData(ctx context.Context, tenantID string, data []byte) (*models.Feedback, error) {
	// Define a structure that matches the expected format of the Intercom JSON data
	var intercomData struct {
		ID             string `json:"id"`
		ConversationID string `json:"conversation_id"`
		Messages       []struct {
			ID     string `json:"id"`
			Author struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"author"`
			Body      string    `json:"body"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"conversation_parts"`
		CreatedAt time.Time `json:"created_at"`
	}

	// Unmarshal the raw JSON data into the structure
	if err := json.Unmarshal(data, &intercomData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Intercom data: %v", err)
	}

	// Convert the Intercom messages to the internal Message format
	messages := make([]models.Message, len(intercomData.Messages))
	for i, msg := range intercomData.Messages {
		messages[i] = models.Message{
			ID:        msg.ID,
			Author:    msg.Author.ID,
			Content:   msg.Body,
			Timestamp: msg.CreatedAt,
		}
	}

	// Create the ConversationContent struct
	content := models.ConversationContent{
		ConversationID: intercomData.ConversationID,
		Messages:       messages,
	}

	// Create and return the FeedbackRecord
	return &models.Feedback{
		ID:        intercomData.ID,
		TenantID:  tenantID,
		Source:    a.GetSourceName(),
		Type:      "Conversation",
		CreatedAt: intercomData.CreatedAt,
		UpdatedAt: time.Now(),
		Metadata:  map[string]interface{}{},
		Content:   content,
	}, nil
}

// GetSourceName returns the name of the source for this adapter
func (a *IntercomIntegration) GetSourceName() string {
	return "Intercom"
}
