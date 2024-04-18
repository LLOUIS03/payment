-- name: ListUsers :many
SELECT  id, name, email, password, create_at, updated_at FROM merchant;
    
-- name: GetUser :one
SELECT  id, name, email, password FROM merchant
WHERE email = $1;