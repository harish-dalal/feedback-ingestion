package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/integrations"
	"github.com/harish-dalal/feedback-ingestion-system/pkg/subscription"
	"github.com/robfig/cron/v3"
)

type CronManager struct {
	subService         *subscription.SubscriptionService
	integrationManager *integrations.IntegrationManager
	cron               *cron.Cron
}

func NewCronManager(subService *subscription.SubscriptionService, integrationManager *integrations.IntegrationManager) *CronManager {
	return &CronManager{
		subService:         subService,
		integrationManager: integrationManager,
		cron:               cron.New(cron.WithSeconds()),
	}
}

func (cm *CronManager) StartGlobalPullJob(ctx context.Context, interval time.Duration) error {
	_, err := cm.cron.AddFunc(fmt.Sprintf("@every %s", interval.String()), func() {
		subscriptions, err := cm.subService.GetAllActivePullSubscriptions(ctx)
		if err != nil {
			fmt.Printf("Failed to query subscriptions: %v\n", err)
			return
		}

		for _, sub := range subscriptions {
			jobCtx, cancel := context.WithTimeout(ctx, interval)
			defer cancel()

			_, err := cm.integrationManager.Pull(jobCtx, sub)
			if err != nil {
				fmt.Printf("Error pulling data for subscription %s: %v\n", sub.ID, err)
			} else {
				err := cm.subService.UpdateLastPulled(jobCtx, sub.ID)
				if err != nil {
					fmt.Printf("Failed to update last pulled time for subscription %s: %v\n", sub.ID, err)
				}
			}
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule global cron job: %v", err)
	}

	cm.cron.Start()
	return nil
}
