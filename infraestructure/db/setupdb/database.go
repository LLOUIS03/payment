package setupdb

import (
	"context"
	"database/sql"
)

type Database interface {
	// Open opens the database connection
	Open() error

	// CloseFunc returns a function that closes the database connection
	CloseFunc() func() error

	// DB returns the database connection
	DB() *sql.DB

	// Ping checks if the database is reachable
	Ping(ctx context.Context) error

	// Setup runs the migrations and sets up the database
	Setup(ctx context.Context) (Database, error)
}
