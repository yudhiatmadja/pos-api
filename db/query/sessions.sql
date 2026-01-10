-- name: CreateTable :one
INSERT INTO tables (outlet_id, name, capacity, qr_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListTables :many
SELECT * FROM tables
WHERE outlet_id = $1
ORDER BY name;

-- name: CreateSession :one
INSERT INTO table_sessions (table_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByToken :one
SELECT ts.*, t.outlet_id, t.name as table_name 
FROM table_sessions ts
JOIN tables t ON ts.table_id = t.id
WHERE ts.token = $1 AND ts.is_active = TRUE AND ts.expires_at > NOW()
LIMIT 1;

-- name: InvalidateSession :exec
UPDATE table_sessions
SET is_active = FALSE
WHERE id = $1;
