package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type SubscriptionHandler struct {
	service *SubscriptionService
}

func NewSubscriptionHandler(service *SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if sub.TenantID == "" || sub.AppID == "" || sub.Source == "" || sub.SubscriptionMode == "" {
		http.Error(w, "TenantID, AppID, SourceID, and SubscriptionMode are required", http.StatusBadRequest)
		return
	}

	if sub.SubscriptionMode != models.SubscriptionModePush && sub.SubscriptionMode != models.SubscriptionModePull {
		http.Error(w, "SubscriptionMode must be 'push' or 'pull'", http.StatusBadRequest)
		return
	}

	sub.ID = uuid.New().String()

	ctx := r.Context()
	if err := h.service.CreateSubscription(ctx, &sub); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create subscription: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}
