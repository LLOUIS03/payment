package setupdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/deuna/payment/utils/retry"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresDB struct {
	db           *sql.DB
	migrations   Migrations
	connStr      string
	maxIdleConns int
	macOpenConns int
	closeFunc    func() error
}

func NewPostgresDB(migrations Migrations, connStr string, maxIdleConns, maxOpenConns int) Database {
	return &postgresDB{
		connStr:      connStr,
		migrations:   migrations,
		maxIdleConns: maxIdleConns,
		macOpenConns: maxOpenConns,
	}
}

func (d *postgresDB) CloseFunc() func() error {
	return d.closeFunc
}

func (d *postgresDB) Open() (err error) {
	d.db, err = sql.Open("pgx", d.connStr)
	if err != nil {
		return err
	}

	d.db.SetMaxIdleConns(d.maxIdleConns)
	d.db.SetMaxOpenConns(d.macOpenConns)

	d.closeFunc = func() error {
		if err := d.db.Close(); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (d *postgresDB) Ping(ctx context.Context) error {
	err := d.db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *postgresDB) DB() *sql.DB {
	return d.db
}

func (d *postgresDB) Setup(ctx context.Context) (Database, error) {
	err := d.Open()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}

	res, err := retry.NewRetry(func() error {
		return d.Ping(ctx)
	}).WithMaxAttempts(3).
		WithSleep(time.Second).
		Run()

	if err != nil {
		msg := ""
		for i := 0; i < res.Attempts(); i++ {
			msg += fmt.Sprintf("Attempt %d: %s\n", i, res[i+1])
		}
		return nil, errors.Wrapf(err, "pinging database\n%s", msg)
	}

	if d.migrations != nil {
		d.migrations.SetDB(d.db)
		err = d.migrations.Run(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "running migrations")
		}
	}

	return d, nil
}
