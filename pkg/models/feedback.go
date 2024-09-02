package models

import "time"

// FeedbackRecord represents the unified structure for feedback
type Feedback struct {
	ID        string                 `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	Source    string                 `json:"source"`
	Type      string                 `json:"type"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata"`
	Content   interface{}            `json:"content"` // Holds type-specific content
}

// Message represents a single message in a conversation
type Message struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationContent represents content specific to conversations
type ConversationContent struct {
	ConversationID string    `json:"conversation_id"`
	Messages       []Message `json:"messages"`
}
