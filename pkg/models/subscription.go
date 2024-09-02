package models

import "time"

type Subscription struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenant_id"`
	AppID            string                 `json:"app_id"`
	SourceID         string                 `json:"source_id"`
	SubscriptionMode string                 `json:"subscriptionMode"` // 'push' or 'pull'
	Configuration    map[string]interface{} `json:"configuration"`
	CreatedAt        time.Time              `json:"created_at"`
	LastPulled       time.Time              `json:"last_pulled"` // Only applicable for pull subscriptions
	Active           bool                   `json:"active"`
}
