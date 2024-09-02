package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

// SubscriptionHandler provides HTTP handlers for managing subscriptions
type SubscriptionHandler struct {
	service *SubscriptionService
}

// NewSubscriptionHandler creates a new SubscriptionHandler
func NewSubscriptionHandler(service *SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscriptionHandler handles HTTP requests to create a new subscription
func (h *SubscriptionHandler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if sub.TenantID == "" || sub.AppID == "" || sub.SourceID == "" || sub.SubscriptionMode == "" {
		http.Error(w, "TenantID, AppID, SourceID, and Model are required", http.StatusBadRequest)
		return
	}

	if sub.SubscriptionMode != "push" && sub.SubscriptionMode != "pull" {
		http.Error(w, "Model must be 'push' or 'pull'", http.StatusBadRequest)
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

// GetSubscriptionHandler handles HTTP requests to retrieve a subscription by ID
// func (h *SubscriptionHandler) GetSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
// 	subscriptionID := r.URL.Query().Get("id")
// 	if subscriptionID == "" {
// 		http.Error(w, "Subscription ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	sub, err := h.service.GetSubscription(ctx, subscriptionID)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to retrieve subscription: %v", err), http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(sub)
// }

// DeleteSubscriptionHandler handles HTTP requests to delete a subscription by ID
// func (h *SubscriptionHandler) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
// 	subscriptionID := r.URL.Query().Get("id")
// 	if subscriptionID == "" {
// 		http.Error(w, "Subscription ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	if err := h.service.DeleteSubscription(ctx, subscriptionID); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to delete subscription: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// ListSubscriptionsByTenantAndAppHandler handles HTTP requests to list all subscriptions for a tenant's app
// func (h *SubscriptionHandler) ListSubscriptionsByTenantAndAppHandler(w http.ResponseWriter, r *http.Request) {
// 	tenantID := r.URL.Query().Get("tenant_id")
// 	appID := r.URL.Query().Get("app_id")
// 	if tenantID == "" || appID == "" {
// 		http.Error(w, "Tenant ID and App ID are required", http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	subs, err := h.service.ListSubscriptionsByTenantAndApp(ctx, tenantID, appID)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to list subscriptions: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(subs)
// }
