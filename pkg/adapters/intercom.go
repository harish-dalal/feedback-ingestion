package adapters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// IntercomAdapter implements the SourceAdapter interface for Intercom data
type IntercomAdapter struct{}

// ProcessRawData processes raw Intercom data and returns a FeedbackRecord
func (a *IntercomAdapter) ProcessRawData(tenantID string, data []byte) (*models.FeedbackRecord, error) {
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
	return &models.FeedbackRecord{
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
func (a *IntercomAdapter) GetSourceName() string {
	return "Intercom"
}
