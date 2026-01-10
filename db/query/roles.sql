-- name: CreateRole :one
INSERT INTO roles (code, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE code = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY name;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_code)
VALUES ($1, $2);

-- name: GetUserRoles :many
SELECT r.code, r.name 
FROM roles r
JOIN user_roles ur ON r.code = ur.role_code
WHERE ur.user_id = $1;
