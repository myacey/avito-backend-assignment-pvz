-- name: CreateReception :one
INSERT INTO receptions (id, date_time, pvz_id) VALUES
($1, $2, $3)
RETURNING *;

-- name: GetOpenReceptionByPvzID :one
SELECT * FROM receptions
WHERE pvz_id = $1 AND status = 'in_progress'
LIMIT 1;

-- name: GetReceptionsByTime :many
SELECT * FROM receptions
WHERE date_time BETWEEN $1 AND $2;

-- name: GetReceptionsByPvzAndTime :many
SELECT * FROM receptions
WHERE pvz_id IN ($1) AND date_time BETWEEN $2 AND $3;

-- name: AddProductToReception :one
INSERT INTO products (id, type, reception_id) VALUES
($1, $2, $3)
RETURNING *;

-- name: GetProductsFromReception :many
SELECT * FROM products
WHERE reception_id IN ($1);

-- name: GetLastProductInReception :one
SELECT * FROM products
WHERE reception_id = $1
ORDER BY date_time DESC
LIMIT 1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: FinishReception :one
UPDATE receptions
SET status='close'
WHERE id=(SELECT id FROM receptions R WHERE R.pvz_id=$1 AND R.status='in_progress' LIMIT 1)
RETURNING *;