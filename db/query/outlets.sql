-- name: CreateOutlet :one
INSERT INTO outlets (name, address, phone)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOutlet :one
SELECT * FROM outlets
WHERE id = $1 LIMIT 1;

-- name: ListOutlets :many
SELECT * FROM outlets
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateOutlet :one
UPDATE outlets
SET name = $2, address = $3, phone = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOutlet :exec
DELETE FROM outlets
WHERE id = $1;
