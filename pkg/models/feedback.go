package models

import "time"

type Feedback struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Source      Source                 `json:"source"`
	SubSourceID string                 `json:"sub_source_id"` // defining either tag / app or relevant identifier
	Type        string                 `json:"type"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata"`
	Content     interface{}            `json:"content"` // Holds type-specific content
}

type Message struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ConversationContent struct {
	ConversationID string    `json:"conversation_id"`
	Messages       []Message `json:"messages"`
}

// Define various other contents depending on the source
