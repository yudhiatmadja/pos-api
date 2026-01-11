-- name: CreateShift :one
INSERT INTO shifts (
    user_id, outlet_id, start_cash
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CloseShift :one
UPDATE shifts
SET end_time = NOW(),
    end_cash = $2,
    expected_cash = $3
WHERE id = $1
RETURNING *;

-- name: GetCurrentShift :one
SELECT * FROM shifts
WHERE user_id = $1 AND end_time IS NULL
LIMIT 1;

-- name: GetShiftById :one
SELECT * FROM shifts
WHERE id = $1 LIMIT 1;

-- name: ListShifts :many
SELECT * FROM shifts
WHERE outlet_id = $1
ORDER BY start_time DESC
LIMIT $2 OFFSET $3;
