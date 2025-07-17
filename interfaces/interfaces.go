package interfaces

import (
	"context"
	"time"
)

// =============================================================================
// CORE INTERFACES (Open Source - Safe for GitHub)
// =============================================================================

// TodoService defines the todo service interface
type TodoService interface {
	CreateTodo(ctx context.Context, todo *Todo) (*Todo, error)
	GetTodo(ctx context.Context, id string) (*Todo, error)
	ListTodos(ctx context.Context, filter TodoFilter) ([]*Todo, error)
	UpdateTodo(ctx context.Context, id string, updates TodoUpdates) (*Todo, error)
	DeleteTodo(ctx context.Context, id string) error
}

// WebhookService handles webhook notifications
type WebhookService interface {
	TriggerWebhook(ctx context.Context, event WebhookEvent) error
	RegisterWebhook(ctx context.Context, webhook *Webhook) error
}

// AnalyticsService handles usage analytics
type AnalyticsService interface {
	TrackEvent(ctx context.Context, event AnalyticsEvent) error
	GetMetrics(ctx context.Context, period string) (*Metrics, error)
}

// =============================================================================
// SHARED DATA MODELS (Safe for public)
// =============================================================================

type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      string    `json:"user_id"`
}

type TodoFilter struct {
	UserID    string `json:"user_id"`
	Completed *bool  `json:"completed"`
	Priority  string `json:"priority"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

type TodoUpdates struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
	Priority    *string `json:"priority"`
}

type WebhookEvent struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    string      `json:"user_id"`
	Data      interface{} `json:"data"`
}

type Webhook struct {
	ID     string   `json:"id"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Secret string   `json:"secret"`
}

type AnalyticsEvent struct {
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

type Metrics struct {
	EventCount  int            `json:"event_count"`
	UniqueUsers int            `json:"unique_users"`
	ByType      map[string]int `json:"by_type"`
}
