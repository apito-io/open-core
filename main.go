package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"open-core/core/handlers"
	"open-core/core/models"
	"open-core/core/services"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services
	todoStore := models.NewInMemoryTodoStore()
	todoService := services.NewCoreTodoService(todoStore)
	webhookService := services.NewCoreWebhookService()
	analyticsService := services.NewCoreAnalyticsService()

	// Initialize handlers
	handler := handlers.NewTodoHandler(todoService, webhookService, analyticsService)

	// Routes
	api := e.Group("/api/v1")
	api.GET("/todos", handler.ListTodos)
	api.POST("/todos", handler.CreateTodo)
	api.GET("/todos/:id", handler.GetTodo)
	api.PUT("/todos/:id", handler.UpdateTodo)
	api.DELETE("/todos/:id", handler.DeleteTodo)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"version": getVersion(),
			"edition": "open-core",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Apito Open Core server on port %s", port)
	log.Printf("✅ CRUD operations, basic webhooks, simple analytics")
	log.Fatal(e.Start(":" + port))
}

func getVersion() string {
	if version := os.Getenv("APITO_VERSION"); version != "" {
		return version
	}
	return "dev"
}
