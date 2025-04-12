-- name: CreatePVZ :one
INSERT INTO pvz (id, registration_date, city) VALUES
($1, $2, $3)
RETURNING *;

-- name: SearchPVZ :many
SELECT * FROM pvz
OFFSET $1 LIMIT $2;