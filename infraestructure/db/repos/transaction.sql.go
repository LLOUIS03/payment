// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transaction.sql

package repos

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO tx (
    id,
    amount,
    tx_type,
    merchant_id
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING id, updated_at
`

type CreateTransactionParams struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	TxType     int64     `json:"tx_type"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type CreateTransactionRow struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (CreateTransactionRow, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.ID,
		arg.Amount,
		arg.TxType,
		arg.MerchantID,
	)
	var i CreateTransactionRow
	err := row.Scan(&i.ID, &i.UpdatedAt)
	return i, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, merchant_id, amount, tx_type, create_at, updated_at FROM tx
WHERE id = $1
`

func (q *Queries) GetTransaction(ctx context.Context, id uuid.UUID) (Tx, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Tx
	err := row.Scan(
		&i.ID,
		&i.MerchantID,
		&i.Amount,
		&i.TxType,
		&i.CreateAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTransactions = `-- name: ListTransactions :many
SELECT id, merchant_id, amount, tx_type, create_at, updated_at FROM tx
`

func (q *Queries) ListTransactions(ctx context.Context) ([]Tx, error) {
	rows, err := q.db.QueryContext(ctx, listTransactions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tx
	for rows.Next() {
		var i Tx
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.Amount,
			&i.TxType,
			&i.CreateAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransactionsByMechant = `-- name: ListTransactionsByMechant :many
SELECT id, merchant_id, amount, tx_type, create_at, updated_at FROM tx
WHERE merchant_id = $1
`

func (q *Queries) ListTransactionsByMechant(ctx context.Context, merchantID uuid.UUID) ([]Tx, error) {
	rows, err := q.db.QueryContext(ctx, listTransactionsByMechant, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tx
	for rows.Next() {
		var i Tx
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.Amount,
			&i.TxType,
			&i.CreateAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransaction = `-- name: UpdateTransaction :exec
UPDATE tx SET amount = $1, tx_type = $2, updated_at = now()
WHERE id = $3 AND updated_at = $4
`

type UpdateTransactionParams struct {
	Amount    float64   `json:"amount"`
	TxType    int64     `json:"tx_type"`
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, updateTransaction,
		arg.Amount,
		arg.TxType,
		arg.ID,
		arg.UpdatedAt,
	)
	return err
}

const updateTransactionType = `-- name: UpdateTransactionType :one
UPDATE tx SET tx_type = $1, updated_at = now()
WHERE id = $2 AND updated_at = $3 RETURNING id, updated_at
`

type UpdateTransactionTypeParams struct {
	TxType    int64     `json:"tx_type"`
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateTransactionTypeRow struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateTransactionType(ctx context.Context, arg UpdateTransactionTypeParams) (UpdateTransactionTypeRow, error) {
	row := q.db.QueryRowContext(ctx, updateTransactionType, arg.TxType, arg.ID, arg.UpdatedAt)
	var i UpdateTransactionTypeRow
	err := row.Scan(&i.ID, &i.UpdatedAt)
	return i, err
}