package feedback

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type FeedbackHandler struct {
	service *FeedbackService
}

func NewFeedbackHandler(service *FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: service}
}

func (h *FeedbackHandler) CreateFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	var feedback models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if feedback.TenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(feedback.TenantID); err != nil {
		http.Error(w, "Invalid Tenant ID format", http.StatusBadRequest)
		return
	}

	feedback.ID = uuid.New().String()

	ctx := r.Context()
	if err := h.service.CreateFeedback(ctx, &feedback); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create feedback record: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feedback)
}

func (h *FeedbackHandler) GetFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	feedbackID := r.URL.Query().Get("id")
	if feedbackID == "" {
		http.Error(w, "Feedback ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	record, err := h.service.GetFeedback(ctx, feedbackID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve feedback record: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(record)
}

func (h *FeedbackHandler) UpdateFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	var feedback models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if feedback.ID == "" {
		http.Error(w, "Feedback ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.UpdateFeedback(ctx, &feedback); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update feedback record: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(feedback)
}

func (h *FeedbackHandler) DeleteFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	feedbackID := r.URL.Query().Get("id")
	if feedbackID == "" {
		http.Error(w, "Feedback ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteFeedback(ctx, feedbackID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete feedback record: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FeedbackHandler) ListFeedbackByTenantHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	feedbacks, err := h.service.ListFeedbackByTenant(ctx, tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list feedback records: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(feedbacks)
}
