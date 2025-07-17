package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"subtree-example/interfaces"
)

// InMemoryTodoStore implements a simple in-memory store for todos
// This is part of the open-source core
type InMemoryTodoStore struct {
	todos map[string]*interfaces.Todo
	mu    sync.RWMutex
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{
		todos: make(map[string]*interfaces.Todo),
	}
}

func (s *InMemoryTodoStore) Create(ctx context.Context, todo *interfaces.Todo) (*interfaces.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if todo.ID == "" {
		todo.ID = uuid.New().String()
	}

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	s.todos[todo.ID] = todo
	return todo, nil
}

func (s *InMemoryTodoStore) GetByID(ctx context.Context, id string) (*interfaces.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, fmt.Errorf("todo with id %s not found", id)
	}

	// Return a copy to avoid race conditions
	todoCopy := *todo
	return &todoCopy, nil
}

func (s *InMemoryTodoStore) List(ctx context.Context, filter interfaces.TodoFilter) ([]*interfaces.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var todos []*interfaces.Todo
	for _, todo := range s.todos {
		if s.matchesFilter(todo, filter) {
			todoCopy := *todo
			todos = append(todos, &todoCopy)
		}
	}

	// Apply pagination
	start := filter.Offset
	if start > len(todos) {
		return []*interfaces.Todo{}, nil
	}

	end := start + filter.Limit
	if filter.Limit == 0 || end > len(todos) {
		end = len(todos)
	}

	return todos[start:end], nil
}

func (s *InMemoryTodoStore) Update(ctx context.Context, id string, updates interfaces.TodoUpdates) (*interfaces.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, fmt.Errorf("todo with id %s not found", id)
	}

	// Apply updates
	if updates.Title != nil {
		todo.Title = *updates.Title
	}
	if updates.Description != nil {
		todo.Description = *updates.Description
	}
	if updates.Completed != nil {
		todo.Completed = *updates.Completed
	}
	if updates.Priority != nil {
		todo.Priority = *updates.Priority
	}

	todo.UpdatedAt = time.Now()

	todoCopy := *todo
	return &todoCopy, nil
}

func (s *InMemoryTodoStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		return fmt.Errorf("todo with id %s not found", id)
	}

	delete(s.todos, id)
	return nil
}

func (s *InMemoryTodoStore) GetAll(ctx context.Context) ([]*interfaces.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var todos []*interfaces.Todo
	for _, todo := range s.todos {
		todoCopy := *todo
		todos = append(todos, &todoCopy)
	}

	return todos, nil
}

func (s *InMemoryTodoStore) matchesFilter(todo *interfaces.Todo, filter interfaces.TodoFilter) bool {
	if filter.UserID != "" && todo.UserID != filter.UserID {
		return false
	}

	if filter.Completed != nil && todo.Completed != *filter.Completed {
		return false
	}

	if filter.Priority != "" && todo.Priority != filter.Priority {
		return false
	}

	return true
} 