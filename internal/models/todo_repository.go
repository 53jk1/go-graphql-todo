package models

import (
	"errors"
	"github.com/dgryski/trifles/uuid"
	"sync"
)

type Todo struct {
	ID          string
	Title       string
	Description string
	IsCompleted bool
}

type TodoRepository struct {
	mu    sync.RWMutex
	items map[string]*Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		items: make(map[string]*Todo),
		mu:    sync.RWMutex{},
	}
}

func (r *TodoRepository) FindAll() ([]*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var todos []*Todo
	for _, todo := range r.items {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) FindByID(id string) (*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.items[id]
	if !ok {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (r *TodoRepository) Create(todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[todo.ID]; ok {
		return errors.New("todo already exists")
	}
	r.items[todo.ID] = todo
	return nil
}

func (r *TodoRepository) Update(todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[todo.ID]; !ok {
		return errors.New("todo not found")
	}
	r.items[todo.ID] = todo
	return nil
}

func (r *TodoRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return errors.New("todo not found")
	}
	delete(r.items, id)
	return nil
}

func NewID() string {
	return uuid.UUIDv4()
}
