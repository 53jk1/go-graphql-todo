package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v4"
	l "github.com/sirupsen/logrus"

	"github.com/53jk1/go-graphql-todo/internal/graphql/generated"
)

var db *sql.DB

type DB struct {
	conn *pgx.Conn
}

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"isCompleted"`
}

type TodoRepository struct {
	db    *sql.DB
	mu    sync.Mutex
	todos []*Todo
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func NewPostgresDB() (*DB, error) {
	var todos []*Todo

	// Build the connection string inside docker-compose
	connString := "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"

	// Build the connection string outside docker-compose
	// connString := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	// check if the database exists
	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", os.Getenv("DB_NAME")).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		l.Info("database exists")
	} else {
		l.Info("database does not exist")
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", os.Getenv("DB_NAME")))
		if err != nil {
			return nil, err
		}
	}

	// check if the table exists
	var tableExists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = $1)", "todos").Scan(&tableExists)
	if err != nil {
		return nil, err
	}
	if tableExists {
		l.Info("table exists")
	} else {
		l.Info("table does not exist")
		_, err = conn.Exec(context.Background(), "CREATE TABLE todos (id UUID PRIMARY KEY, title TEXT, description TEXT, is_completed BOOLEAN)")
		if err != nil {
			return nil, err
		}
	}

	// load all the todos from the database for the in-memory cache
	rows, err := conn.Query(context.Background(), "SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// create the in-memory cache
	cache := NewTodoRepository(db)
	for _, todo := range todos {
		cache.Add(todo)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) GetConn() *pgx.Conn {
	return db.conn
}

func (db *DB) Close() {
	db.conn.Close(context.Background())
}

func ConnectDB() (*sql.DB, error) {
	// Read PostgreSQL database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	l.Info("Connecting to database")
	l.Info("Host: ", host)
	l.Info("Port: ", port)
	l.Info("User: ", user)
	l.Info("Password: ", password)
	l.Info("DB Name: ", dbname)

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	l.Info("Connection string: ", connStr)

	// Establish connection to PostgreSQL database
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		l.Error("Error connecting to database: ", err)
		return nil, err
	}

	// Set maximum number of connections to the database
	db.SetMaxOpenConns(5)

	// Set maximum number of connections in the idle connection pool
	db.SetMaxIdleConns(5)

	// Set maximum amount of time a connection may be reused
	db.SetConnMaxLifetime(0)

	// Set maximum amount of time to wait for a connection to be returned to the pool
	db.SetConnMaxIdleTime(0)

	// Verify connection to PostgreSQL database
	err = db.Ping()
	if err != nil {
		l.Error("Error connecting to database: ", err)
		return nil, err
	}

	l.Info("Successfully connected to database")

	return db, nil
}

func GetTodos() ([]*Todo, error) {
	// Execute SELECT statement on the todos table
	rows, err := db.Query("SELECT id, title, description, is_completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse rows into []*Todo slice
	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func CreateTodo(todo *Todo) error {

	db, err := ConnectDB()
	if err != nil {
		return err
	}

	// Get incoming data
	id := todo.ID
	title := todo.Title
	description := todo.Description
	isCompleted := todo.IsCompleted

	// Execute INSERT statement on the todos table
	_, err = db.Exec("INSERT INTO todos (id, title, description, is_completed) VALUES ($1, $2, $3, $4)",
		id, title, description, isCompleted)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTodo(todo *Todo) error {
	// Get incoming data
	id := todo.ID
	title := todo.Title
	description := todo.Description
	isCompleted := todo.IsCompleted

	// Execute UPDATE statement on the todos table
	_, err := db.Exec("UPDATE todos SET title=$1, description=$2, is_completed=$3 WHERE id=$4",
		title, description, isCompleted, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTodo(id string) error {
	// Execute DELETE statement on the todos table
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TodoRepository) GetAll() ([]*Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query("SELECT id, title, description, is_completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Add(todo *Todo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.todos = append(r.todos, todo)
}

// GetAllTodos returns all todos
func GetAllTodos() ([]*generated.Todo, error) {
	// Execute SELECT statement on the todos table
	rows, err := db.Query("SELECT id, title, description, is_completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	// Parse rows into []*Todo slice
	var todos []*generated.Todo
	for rows.Next() {
		todo := &generated.Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// DeleteTodoByID deletes a todo by ID
func DeleteTodoByID(id string) error {
	// Execute DELETE statement on the todos table
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

// GetTodoByID returns a todo by ID
func GetTodoByID(id string) (*generated.Todo, error) {
	// Connect to database
	db, err := ConnectDB()
	if err != nil {
		l.Error("Error connecting to database: ", err)
		return nil, err
	}

	// Execute SELECT statement on the todos table
	row := db.QueryRow("SELECT id, title, description, is_completed FROM todos WHERE id=$1", id)

	// Parse row into *Todo
	todo := &generated.Todo{}
	err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// UpdateTodoByID updates a todo by ID
func UpdateTodoByID(id string, todo *generated.NewTodo) error {
	// Execute UPDATE statement on the todos table
	_, err := db.Exec("UPDATE todos SET title=$1, description=$2, is_completed=$3 WHERE id=$4",
		todo.Title, todo.Description, todo.IsCompleted, id)
	if err != nil {
		return err
	}

	return nil
}
