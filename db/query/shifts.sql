-- name: CreateShift :one
INSERT INTO shifts (
    user_id, store_id, opening_cash
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CloseShift :one
UPDATE shifts
SET closed_at = NOW(),
    closing_cash = $2,
    expected_cash = $3
WHERE id = $1
RETURNING *;

-- name: GetCurrentShift :one
SELECT * FROM shifts
WHERE user_id = $1 AND closed_at IS NULL
LIMIT 1;

-- name: ListShifts :many
SELECT * FROM shifts
WHERE store_id = $1
ORDER BY opened_at DESC
LIMIT $2 OFFSET $3;
