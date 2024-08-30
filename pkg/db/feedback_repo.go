package db

import (
	"fmt"
	"sync"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// Mock database for storing feedback records
var feedbackDB = struct {
	sync.RWMutex
	feedbacks map[string]*models.FeedbackRecord
}{
	feedbacks: make(map[string]*models.FeedbackRecord),
}

// FeedbackRepository defines methods for managing feedback in the database
type FeedbackRepository struct{}

// NewFeedbackRepository creates a new FeedbackRepository
func NewFeedbackRepository() *FeedbackRepository {
	return &FeedbackRepository{}
}

// Save saves a new feedback record to the database
func (repo *FeedbackRepository) Save(record *models.FeedbackRecord) error {
	feedbackDB.Lock()
	defer feedbackDB.Unlock()

	if _, exists := feedbackDB.feedbacks[record.ID]; exists {
		return fmt.Errorf("feedback with ID %s already exists", record.ID)
	}

	feedbackDB.feedbacks[record.ID] = record
	return nil
}

// Get retrieves a feedback record by ID from the database
func (repo *FeedbackRepository) Get(feedbackID string) (*models.FeedbackRecord, error) {
	feedbackDB.RLock()
	defer feedbackDB.RUnlock()

	record, exists := feedbackDB.feedbacks[feedbackID]
	if !exists {
		return nil, fmt.Errorf("feedback with ID %s not found", feedbackID)
	}

	return record, nil
}

// Update updates an existing feedback record in the database
func (repo *FeedbackRepository) Update(record *models.FeedbackRecord) error {
	feedbackDB.Lock()
	defer feedbackDB.Unlock()

	if _, exists := feedbackDB.feedbacks[record.ID]; !exists {
		return fmt.Errorf("feedback with ID %s not found", record.ID)
	}

	feedbackDB.feedbacks[record.ID] = record
	return nil
}

// Delete deletes a feedback record by ID from the database
func (repo *FeedbackRepository) Delete(feedbackID string) error {
	feedbackDB.Lock()
	defer feedbackDB.Unlock()

	if _, exists := feedbackDB.feedbacks[feedbackID]; !exists {
		return fmt.Errorf("feedback with ID %s not found", feedbackID)
	}

	delete(feedbackDB.feedbacks, feedbackID)
	return nil
}
