package services

import (
	"context"
	"fmt"

	"open-core/interfaces"
	"open-core/interfaces"
)

// CoreTodoService implements the core TodoService interface
// This is part of the open-source core
type CoreTodoService struct {
	store *models.InMemoryTodoStore
}

func NewCoreTodoService(store *models.InMemoryTodoStore) interfaces.TodoService {
	return &CoreTodoService{
		store: store,
	}
}

func (s *CoreTodoService) CreateTodo(ctx context.Context, todo *interfaces.Todo) (*interfaces.Todo, error) {
	// Basic validation
	if todo.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return s.store.Create(ctx, todo)
}

func (s *CoreTodoService) GetTodo(ctx context.Context, id string) (*interfaces.Todo, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	return s.store.GetByID(ctx, id)
}

func (s *CoreTodoService) ListTodos(ctx context.Context, filter interfaces.TodoFilter) ([]*interfaces.Todo, error) {
	// Set default limit for open source version
	if filter.Limit == 0 {
		filter.Limit = 50 // Basic limit
	}

	return s.store.List(ctx, filter)
}

func (s *CoreTodoService) UpdateTodo(ctx context.Context, id string, updates interfaces.TodoUpdates) (*interfaces.Todo, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	return s.store.Update(ctx, id, updates)
}

func (s *CoreTodoService) DeleteTodo(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return s.store.Delete(ctx, id)
}
