package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

const (
	POSTGRES_USER     = "postgres"
	POSTGRES_PASSWORD = "postgres"
	POSTGRES_DB       = "todo"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB() (*DB, error) {
	// Build the connection string
	connString := "postgres://postgres:postgres@db/todo?sslmode=disable"

	// Create a new connection pool
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	// add username and password to connConfig
	connConfig.User = POSTGRES_USER
	connConfig.Password = POSTGRES_PASSWORD
	connConfig.Database = POSTGRES_DB

	connPool, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Return the DB object
	return &DB{conn: connPool}, nil
}

func (db *DB) Close() {
	db.conn.Close(context.Background())
}
