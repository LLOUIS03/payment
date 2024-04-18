package multiinstruction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/deuna/payment/infraestructure/db/repos"
)

type MultiInstruction interface {
	ExecTx(context.Context, func(repos.Querier) error) error
}

type multiInstruction struct {
	db *sql.DB
}

func NewMultiInstruction(db *sql.DB) MultiInstruction {
	return &multiInstruction{
		db: db,
	}
}

// ExecTx executes a transaction
func (m *multiInstruction) ExecTx(ctx context.Context, fn func(repos.Querier) error) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := repos.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
