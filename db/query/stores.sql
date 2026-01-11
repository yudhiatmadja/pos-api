-- name: CreateStore :one
INSERT INTO stores (
    name, address, phone
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetStore :one
SELECT * FROM stores
WHERE id = $1 LIMIT 1;

-- name: ListStores :many
SELECT * FROM stores
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateStore :one
UPDATE stores
SET name = $2, address = $3, phone = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteStore :exec
DELETE FROM stores
WHERE id = $1;
