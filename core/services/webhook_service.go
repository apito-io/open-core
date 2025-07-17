package services

import (
	"context"
	"fmt"
	"log"

	"open-core/interfaces"
)

// CoreWebhookService provides basic webhook functionality for open source
type CoreWebhookService struct{}

func NewCoreWebhookService() interfaces.WebhookService {
	return &CoreWebhookService{}
}

func (s *CoreWebhookService) TriggerWebhook(ctx context.Context, event interfaces.WebhookEvent) error {
	// Basic webhook - just log the event
	log.Printf("Webhook triggered: %s for user %s", event.Type, event.UserID)
	return nil
}

func (s *CoreWebhookService) RegisterWebhook(ctx context.Context, webhook *interfaces.Webhook) error {
	// Basic implementation - just validate
	if webhook.URL == "" {
		return fmt.Errorf("webhook URL is required")
	}
	log.Printf("Webhook registered: %s", webhook.URL)
	return nil
}
