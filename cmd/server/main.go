package main

import (
	"fmt"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/adapters"
)

func main() {
	// Example raw JSON data from Intercom
	intercomData := []byte(`
    {
        "id": "12345",
        "conversation_id": "conv_123",
        "created_at": "2023-08-30T12:34:56Z",
        "conversation_parts": [
            {
                "id": "msg_1",
                "author": {"id": "user_1", "type": "user"},
                "body": "Hello, I need help!",
                "created_at": "2023-08-30T12:34:56Z"
            },
            {
                "id": "msg_2",
                "author": {"id": "agent_1", "type": "admin"},
                "body": "Sure, how can I assist you?",
                "created_at": "2023-08-30T12:35:30Z"
            }
        ]
    }`)

	// Initialize the IntercomAdapter
	intercomAdapter := &adapters.IntercomAdapter{}

	// Process the data using the adapter
	feedbackRecord, err := intercomAdapter.ProcessRawData("tenant1", intercomData)
	if err != nil {
		fmt.Printf("Error processing Intercom data: %v\n", err)
		return
	}

	// Print the processed feedback record
	fmt.Printf("Processed Feedback Record: %+v\n", feedbackRecord)
}
