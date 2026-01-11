-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    user_id, action, entity, entity_id, "before", "after", ip_address
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
