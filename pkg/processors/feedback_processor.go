package processor

import (
	"fmt"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/adapters"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/tenant"
)

// FeedbackProcessor processes feedback data using a source adapter and tenant information
type FeedbackProcessor struct {
	adapter     adapters.SourceAdapter
	tenanManger *tenant.TenantManager
}

// NewFeedbackProcessor creates a new FeedbackProcessor with the given source adapter
func NewFeedbackProcessor(adapter adapters.SourceAdapter, tenantManager *tenant.TenantManager) *FeedbackProcessor {
	return &FeedbackProcessor{
		adapter:     adapter,
		tenanManger: tenantManager,
	}
}

// ProcessFeedback processes raw feedback data using the configured adapter and tenant ID
func (p *FeedbackProcessor) ProcessFeedback(tenantID string, rawData []byte) (*models.FeedbackRecord, error) {
	tenant, err := p.tenanManger.GetTenant(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}

	// Check if the source is active for this tenant
	isSourceActive := false
	for _, source := range tenant.ActiveSources {
		if source == p.adapter.GetSourceName() { // You'd need to add this method to the SourceAdapter interface
			isSourceActive = true
			break
		}
	}
	if !isSourceActive {
		return nil, fmt.Errorf("source is not active for this tenant")
	}

	return p.adapter.ProcessRawData(tenantID, rawData)
}
