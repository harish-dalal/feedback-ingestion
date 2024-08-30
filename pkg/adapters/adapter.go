package adapters

import "github.com/harish-dalal/feedback-ingestion-system/pkg/models"

// SourceAdapter defines the contract for processing raw feedback data from different sources
type SourceAdapter interface {
	// ProcessRawData processes the raw data and returns a FeedbackRecord
	ProcessRawData(tenantID string, data []byte) (*models.FeedbackRecord, error)

	// GetSourceName returns the name of the source (e.g., "Intercom", "PlayStore")
	GetSourceName() string
}
