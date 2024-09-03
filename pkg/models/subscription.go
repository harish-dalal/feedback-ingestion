package models

import "time"

type Subscription struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenant_id"`
	SubSourceId      string                 `json:"sub_source_id"`
	Source           Source                 `json:"source"` // can be converted to a separate source table for not depending on the source name
	SubscriptionMode SubscriptionMode       `json:"subscription_mode"`
	Configuration    map[string]interface{} `json:"configuration"`
	CreatedAt        time.Time              `json:"created_at"`
	LastPulled       time.Time              `json:"last_pulled"` // Only applicable for pull
	Active           bool                   `json:"active"`
}

type SubscriptionMode string

const (
	SubscriptionModePush SubscriptionMode = "push"
	SubscriptionModePull SubscriptionMode = "pull"
)
