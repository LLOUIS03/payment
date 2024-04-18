package setupdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pressly/goose"
)

type Migrations interface {
	// Run runs the migrations
	Run(ctx context.Context) error

	// SetDB sets the database connection
	SetDB(db *sql.DB)
}

type gooseMigrations struct {
	sourceFolder string
	db           *sql.DB
	logger       goose.Logger
}

func NewGooseMigrations(sourceFolder string, logger goose.Logger) Migrations {
	return &gooseMigrations{
		sourceFolder: sourceFolder,
		logger:       logger,
	}
}

func (g *gooseMigrations) SetDB(db *sql.DB) {
	g.db = db
}

func (g *gooseMigrations) Run(ctx context.Context) error {
	for {
		_, err := g.db.ExecContext(ctx, "CREATE TABLE goose_migrations_in_progress (dummy boolean)")
		if err == nil {
			break
		}

		log.Errorf("Error creating goose_migrations_in_progress table: %v", err)
		time.Sleep(time.Second)
	}

	defer func() {
		_, err := g.db.ExecContext(ctx, "DROP TABLE goose_migrations_in_progress")
		if err != nil {
			log.Errorf("Error dropping goose_migrations_in_progress table: %v", err)
		}
	}()

	goose.SetLogger(g.logger)

	return goose.Run("up", g.db, g.sourceFolder)
}
