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
	strategiesMap := map[models.Source]integrations.SourceStrategy{
		models.SourceIntercom:  integrations.NewIntercomStrategy(),
		models.SourceDiscourse: integrations.NewDiscourseStrategy(),
	}

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

	integrationManager := integrations.NewIntegrationManager(strategiesMap, feedbackService)

	// Init cron manager
	cronManager := cron.NewCronManager(subService, integrationManager)
	// TODO: need to change timer to 8 hr
	err := cronManager.StartGlobalPullJob(context.Background(), 8*time.Hour)
	if err != nil {
		log.Fatalf("Failed to start global pull job: %v", err)
	}

	// webhooks - need to setup web hook routes for all the sources which can support push based ingestion
	srv.Router.HandleFunc("/webhook/discourse", func(w http.ResponseWriter, r *http.Request) {
		integrationManager.HandleWebhook(w, r, models.SourceDiscourse)
	})
	srv.Router.HandleFunc("/webhook/intercom", func(w http.ResponseWriter, r *http.Request) {
		integrationManager.HandleWebhook(w, r, models.SourceIntercom)
	})

	// Health check
	srv.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Tenant CRUD routes
	srv.Router.HandleFunc("/tenant", tenantHandler.CreateTenantHandler)
	srv.Router.HandleFunc("/tenant/get", tenantHandler.GetTenantHandler)
	srv.Router.HandleFunc("/tenant/update", tenantHandler.UpdateTenantHandler)
	srv.Router.HandleFunc("/tenant/delete", tenantHandler.DeleteTenantHandler)

	// Feedback CRUD routes
	srv.Router.HandleFunc("/feedback", feedbackHandler.CreateFeedbackHandler)
	srv.Router.HandleFunc("/feedback/get", feedbackHandler.GetFeedbackHandler)
	srv.Router.HandleFunc("/feedback/update", feedbackHandler.UpdateFeedbackHandler)
	srv.Router.HandleFunc("/feedback/delete", feedbackHandler.DeleteFeedbackHandler)
	srv.Router.HandleFunc("/feedback/list", feedbackHandler.ListFeedbackByTenantHandler)

	// Subscription CRUD routes
	srv.Router.HandleFunc("/subscription", subHandler.CreateSubscriptionHandler)
}
