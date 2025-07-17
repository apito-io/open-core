package services

import (
	"context"
	"log"

	"open-core/interfaces"
)

// CoreAnalyticsService provides basic analytics for open source
type CoreAnalyticsService struct{}

func NewCoreAnalyticsService() interfaces.AnalyticsService {
	return &CoreAnalyticsService{}
}

func (s *CoreAnalyticsService) TrackEvent(ctx context.Context, event interfaces.AnalyticsEvent) error {
	// Basic tracking - just log
	log.Printf("Event tracked: %s for user %s", event.Type, event.UserID)
	return nil
}

func (s *CoreAnalyticsService) GetMetrics(ctx context.Context, period string) (*interfaces.Metrics, error) {
	// Basic metrics - return minimal data
	return &interfaces.Metrics{
		EventCount:  0,
		UniqueUsers: 0,
		ByType:      make(map[string]int),
	}, nil
} 