package integrations

import (
	"fmt"
	"io"
	"net/http"
	"time"

	processor "github.com/harish-dalal/feedback-ingestion-system/pkg/processors"
)

// PullIntegration periodically fetches data from a source
type PullIntegration struct {
	processor *processor.FeedbackProcessor
	interval  time.Duration
	fetchURL  string
}

func NewPullIntegration(processor *processor.FeedbackProcessor, interval time.Duration, fetchURL string) *PullIntegration {
	return &PullIntegration{
		processor: processor,
		interval:  interval,
		fetchURL:  fetchURL,
	}
}

func (pi *PullIntegration) Start() {
	ticker := time.NewTicker(pi.interval)
	go func() {
		for range ticker.C {
			resp, err := http.Get(pi.fetchURL)
			if err != nil {
				fmt.Printf("Failed to fetch data: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to read response body: %v\n", err)
				continue
			}

			tenantID := "tenant1" // In a real implementation, you would map the API key to a tenant ID
			_, err = pi.processor.ProcessFeedback(tenantID, data)
			if err != nil {
				fmt.Printf("Failed to process feedback: %v\n", err)
				continue
			}
		}
	}()
}
