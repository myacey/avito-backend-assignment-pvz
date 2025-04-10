-- name: CreatePVZ :one
INSERT INTO pvz (id, registration_date, city) VALUES
($1, $2, $3)
RETURNING *;

-- name: SearchPvz :many
SELECT * FROM pvz
WHERE id IN ($1)
ORDER BY registration_date DESC
LIMIT $2 OFFSET $3;