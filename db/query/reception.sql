-- name: CreateReception :one
INSERT INTO receptions (id, date_time, pvz_id) VALUES
($1, $2, $3)
RETURNING *;

-- name: GetOpenReceptionByPvz :one
SELECT id FROM receptions
WHERE pvz_id = $1 AND status = 'in_progress'
LIMIT 1;

-- name: GetReceptionsByTime :many
SELECT id, date_time, pvz_id, status
FROM receptions
WHERE date_time BETWEEN $1 AND $2;

-- name: GetReceptionsByPvzAndTime :many
SELECT * FROM receptions
WHERE pvz_id IN ($1) AND date_time BETWEEN $2 AND $3;

-- name: AddProductToReception :one
INSERT INTO products (id, date_time, type, reception_id) VALUES
($1, $2, $3, $4)
RETURNING *;

-- name: GetProductsFromReception :many
SELECT * FROM products
WHERE reception_id IN ($1);

-- name: DeleteProductFromReception :exec
DELETE FROM products
WHERE id = $1;

-- name: FinishReception :exec
UPDATE receptions
SET status='finished'
WHERE id=$1;