package graphql

import (
	"context"

	"github.com/53jk1/go-graphql-todo/internal/graphql/generated"
	"github.com/53jk1/go-graphql-todo/internal/models"
	_ "github.com/99designs/gqlgen/graphql"
)

type Resolver struct {
	TodoRepo *models.TodoRepository
}

func NewResolver(todoRepo *models.TodoRepository) *Resolver {
	return &Resolver{todoRepo}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func NewTodoFromModel(todo *models.Todo) *generated.Todo {
	return &generated.Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}
}

func NewTodoListFromModel(todos []*models.Todo) []*generated.Todo {
	var todoList []*generated.Todo
	for _, todo := range todos {
		todoList = append(todoList, NewTodoFromModel(todo))
	}
	return todoList
}

func (r *mutationResolver) CreateTodo(ctx context.Context, title string, description string) (*generated.Todo, error) {
	todo := &models.Todo{
		ID:          models.NewID(),
		Title:       title,
		Description: description,
		IsCompleted: false,
	}
	if err := r.TodoRepo.Create(todo); err != nil {
		return nil, err
	}
	return NewTodoFromModel(todo), nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, title *string, description *string, isCompleted *bool) (*generated.Todo, error) {
	todo, err := r.TodoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if title != nil {
		todo.Title = *title
	}
	if description != nil {
		todo.Description = *description
	}
	if isCompleted != nil {
		todo.IsCompleted = *isCompleted
	}
	if err := r.TodoRepo.Update(todo); err != nil {
		return nil, err
	}
	return NewTodoFromModel(todo), nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	err := r.TodoRepo.Delete(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*generated.Todo, error) {
	todos, err := r.TodoRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return NewTodoListFromModel(todos), nil
}

func (r *queryResolver) TodoByID(ctx context.Context, id string) (*generated.Todo, error) {
	todo, err := r.TodoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return NewTodoFromModel(todo), nil
}

type QueryResolver interface {
	Hello(ctx context.Context, name string) (string, error)
}
