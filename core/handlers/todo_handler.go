package handlers

import (
	"net/http"
	"strconv"
	"time"

	"open-core/interfaces"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	todoService      interfaces.TodoService
	webhookService   interfaces.WebhookService
	analyticsService interfaces.AnalyticsService
}

func NewTodoHandler(todoService interfaces.TodoService, webhookService interfaces.WebhookService, analyticsService interfaces.AnalyticsService) *TodoHandler {
	return &TodoHandler{
		todoService:      todoService,
		webhookService:   webhookService,
		analyticsService: analyticsService,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var todo interfaces.Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Set user ID from context (in real app, this would come from JWT)
	todo.UserID = c.Request().Header.Get("X-User-ID")
	if todo.UserID == "" {
		todo.UserID = "anonymous"
	}

	created, err := h.todoService.CreateTodo(c.Request().Context(), &todo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Track analytics event
	h.analyticsService.TrackEvent(c.Request().Context(), interfaces.AnalyticsEvent{
		Type:      "todo_created",
		UserID:    created.UserID,
		Data:      map[string]interface{}{"todo_id": created.ID},
		Timestamp: time.Now(),
	})

	// Trigger webhook
	h.webhookService.TriggerWebhook(c.Request().Context(), interfaces.WebhookEvent{
		Type:      "todo.created",
		Timestamp: time.Now(),
		UserID:    created.UserID,
		Data:      created,
	})

	return c.JSON(http.StatusCreated, created)
}

func (h *TodoHandler) GetTodo(c echo.Context) error {
	id := c.Param("id")

	todo, err := h.todoService.GetTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) ListTodos(c echo.Context) error {
	filter := interfaces.TodoFilter{
		UserID: c.Request().Header.Get("X-User-ID"),
	}

	if completed := c.QueryParam("completed"); completed != "" {
		if comp, err := strconv.ParseBool(completed); err == nil {
			filter.Completed = &comp
		}
	}

	if priority := c.QueryParam("priority"); priority != "" {
		filter.Priority = priority
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	if offset := c.QueryParam("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	todos, err := h.todoService.ListTodos(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  todos,
		"count": len(todos),
	})
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id := c.Param("id")

	var updates interfaces.TodoUpdates
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updated, err := h.todoService.UpdateTodo(c.Request().Context(), id, updates)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Track analytics event
	h.analyticsService.TrackEvent(c.Request().Context(), interfaces.AnalyticsEvent{
		Type:      "todo_updated",
		UserID:    updated.UserID,
		Data:      map[string]interface{}{"todo_id": updated.ID},
		Timestamp: time.Now(),
	})

	// Trigger webhook
	h.webhookService.TriggerWebhook(c.Request().Context(), interfaces.WebhookEvent{
		Type:      "todo.updated",
		Timestamp: time.Now(),
		UserID:    updated.UserID,
		Data:      updated,
	})

	return c.JSON(http.StatusOK, updated)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id := c.Param("id")
	userID := c.Request().Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	err := h.todoService.DeleteTodo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Track analytics event
	h.analyticsService.TrackEvent(c.Request().Context(), interfaces.AnalyticsEvent{
		Type:      "todo_deleted",
		UserID:    userID,
		Data:      map[string]interface{}{"todo_id": id},
		Timestamp: time.Now(),
	})

	// Trigger webhook
	h.webhookService.TriggerWebhook(c.Request().Context(), interfaces.WebhookEvent{
		Type:      "todo.deleted",
		Timestamp: time.Now(),
		UserID:    userID,
		Data:      map[string]interface{}{"id": id},
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Todo deleted successfully"})
}
