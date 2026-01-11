-- name: CreateSession :one
INSERT INTO table_sessions (
    table_id, token, expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM table_sessions
WHERE token = $1 LIMIT 1;

-- name: GetTableSessions :many
SELECT ts.*, t.store_id, t.name as table_name 
FROM table_sessions ts
JOIN tables t ON ts.table_id = t.id
WHERE t.store_id = $1
ORDER BY ts.created_at DESC;
