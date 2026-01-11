-- name: CreateIdempotencyKey :one
INSERT INTO idempotency_keys (
    key, response_status, response_body
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetIdempotencyKey :one
SELECT * FROM idempotency_keys
WHERE key = $1 LIMIT 1;
