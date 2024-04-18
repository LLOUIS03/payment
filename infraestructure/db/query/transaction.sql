-- name: CreateTransaction :one
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
) RETURNING id, updated_at;

-- name: ListTransactions :many
SELECT id, merchant_id, amount, tx_type, create_at, updated_at FROM tx;

-- name: ListTransactionsByMechant :many
SELECT id, merchant_id, amount, tx_type, create_at, updated_at FROM tx
WHERE merchant_id = $1;

-- name: GetTransaction :one
SELECT id, merchant_id, amount, tx_type, create_at, updated_at FROM tx
WHERE id = $1;

-- name: UpdateTransaction :exec
UPDATE tx SET amount = $1, tx_type = $2, updated_at = now()
WHERE id = $3 AND updated_at = $4;

-- name: UpdateTransactionType :one
UPDATE tx SET tx_type = $1, updated_at = now()
WHERE id = $2 AND updated_at = $3 RETURNING id, updated_at;