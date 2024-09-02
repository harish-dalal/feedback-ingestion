package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/cron"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/db"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/feedback"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/integrations"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/server"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/subscription"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/tenant"
)

func SetupRoutes(srv *server.Server) {
	// Initialize IntegrationManager with strategies
	strategiesMap := map[models.SourceType]integrations.SourceStrategy{
		models.STIntercom:  integrations.NewIntercomStrategy(),
		models.STDiscourse: integrations.NewDiscourseStrategy(),
	}
	integrationManager := integrations.NewIntegrationManager(strategiesMap)

	// Tenant handlers
	tenantRepo := db.NewTenantRepository(srv.DBPool)
	tenantService := tenant.NewTenantService(tenantRepo)
	tenantHandler := tenant.NewTenantHandler(tenantService)

	// Feedback handlers
	feedbackRepo := db.NewFeedbackRepository(srv.DBPool)
	feedbackService := feedback.NewFeedbackService(feedbackRepo)
	feedbackHandler := feedback.NewFeedbackHandler(feedbackService)

	// Subsription handlers
	subRepo := db.NewSubscriptionRepository(srv.DBPool)
	subService := subscription.NewSubscriptionService(subRepo)
	subHandler := subscription.NewSubscriptionHandler(subService)

	// Init cron manager
	cronManager := cron.NewCronManager(subService, integrationManager)
	// TODO: need to change timer to 8 hr
	err := cronManager.StartGlobalPullJob(context.Background(), 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to start global pull job: %v", err)
	}

	// Set up routes
	srv.Router.HandleFunc("/webhook/discourse", func(w http.ResponseWriter, r *http.Request) {
		integrationManager.HandleWebhook(w, r, models.STDiscourse)
	})
	srv.Router.HandleFunc("/webhook/intercom", func(w http.ResponseWriter, r *http.Request) {
		integrationManager.HandleWebhook(w, r, models.STIntercom)
	})

	// Health check
	srv.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Tenant CRUD routes
	srv.Router.HandleFunc("/tenants", tenantHandler.CreateTenantHandler)
	srv.Router.HandleFunc("/tenants/get", tenantHandler.GetTenantHandler)
	srv.Router.HandleFunc("/tenants/update", tenantHandler.UpdateTenantHandler)
	srv.Router.HandleFunc("/tenants/delete", tenantHandler.DeleteTenantHandler)

	// Feedback CRUD routes
	srv.Router.HandleFunc("/feedback", feedbackHandler.CreateFeedbackHandler)
	srv.Router.HandleFunc("/feedback/get", feedbackHandler.GetFeedbackHandler)
	srv.Router.HandleFunc("/feedback/update", feedbackHandler.UpdateFeedbackHandler)
	srv.Router.HandleFunc("/feedback/delete", feedbackHandler.DeleteFeedbackHandler)
	srv.Router.HandleFunc("/feedback/list", feedbackHandler.ListFeedbackByTenantHandler)

	// Subscription CRUD routes
	srv.Router.HandleFunc("/subscriptions", subHandler.CreateSubscriptionHandler)
}
