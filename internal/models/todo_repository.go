package models

import (
	"errors"
	"sync"

	"github.com/dgryski/trifles/uuid"
	l "github.com/sirupsen/logrus"

	"github.com/53jk1/go-graphql-todo/internal/database"
	"github.com/53jk1/go-graphql-todo/internal/graphql/generated"
)

type TodoRepository struct {
	mu    sync.RWMutex
	items map[string]*generated.Todo
	db    *database.DB
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		items: make(map[string]*generated.Todo),
		mu:    sync.RWMutex{},
		db:    &database.DB{},
	}
}

func (r *TodoRepository) FindAll() ([]*generated.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var todos []*generated.Todo
	for _, todo := range r.items {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) FindByID(id string) (*generated.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if todo, ok := r.items[id]; ok {
		return todo, nil
	}
	return nil, errors.New("todo not found")
}

func (r *TodoRepository) Create(todo *generated.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[todo.ID]; ok {
		return errors.New("todo already exists")
	}
	r.items[todo.ID] = &generated.Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}
	return nil
}

func (r *TodoRepository) Update(id string, todo *generated.NewTodo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		l.Error("todo not found")
		return errors.New("todo not found")
	}
	r.items[id] = &generated.Todo{
		ID:          id,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}
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
