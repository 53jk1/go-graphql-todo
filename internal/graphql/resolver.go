package graphql

import (
	"context"
	"database/sql"
	"log"

	"github.com/dgryski/trifles/uuid"
	l "github.com/sirupsen/logrus"

	"github.com/53jk1/go-graphql-todo/internal/database"

	_ "github.com/99designs/gqlgen/graphql"

	"github.com/53jk1/go-graphql-todo/internal/graphql/generated"
	"github.com/53jk1/go-graphql-todo/internal/models"
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

func (r *mutationResolver) CreateTodo(ctx context.Context, input generated.NewTodo) (*generated.Todo, error) {
	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	todo := &models.Todo{
		ID:          uuid.UUIDv4(),
		Title:       input.Title,
		Description: input.Description,
		IsCompleted: false,
	}
	if err := r.TodoRepo.Create((*generated.Todo)(todo)); err != nil {
		return nil, err
	}

	// Insert todo into database
	if err := database.CreateTodo((*database.Todo)(todo)); err != nil {
		return nil, err
	}

	return (*generated.Todo)(todo), nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, todo generated.NewTodo) (*generated.Todo, error) {
	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	if _, err := database.GetTodoByID(id); err != nil {
		l.Println("GetTodoByID: ", err)
		return nil, err
	}

	// Update todo in database
	if err := database.UpdateTodoByID(id, &todo); err != nil {
		l.Println("UpdateTodoByID: ", err)
		return nil, err
	}

	// map NewTodo to Todo
	t := &models.Todo{
		ID:          id,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}

	return (*generated.Todo)(t), nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (*generated.DeleteTodoPayload, error) {
	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		log.Println("ConnectDB: ", err)
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	if _, err := database.GetTodoByID(id); err != nil {
		log.Println("GetTodoByID: ", err)
		return nil, err
	}

	// Delete todo from database
	if err := database.DeleteTodoByID(id); err != nil {
		log.Println("DeleteTodoByID: ", err)
		return nil, err
	}

	return &generated.DeleteTodoPayload{Message: "Todo deleted successfully"}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*generated.Todo, error) {
	todos, err := r.TodoRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

	// Get all todos from database
	todos, err = database.GetAllTodos()
	if err != nil {
		return nil, err
	}

	// Close database connection
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	return todos, nil
}

func (r *queryResolver) TodoByID(ctx context.Context, id string) (*generated.Todo, error) {
	todo, err := r.TodoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Connect to database
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

	// Get todo by id from database
	todo, err = database.GetTodoByID(id)
	if err != nil {
		return nil, err
	}

	// Close database connection
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	return (*generated.Todo)(todo), nil
}

type QueryResolver interface {
	Hello(ctx context.Context, name string) (string, error)
}
