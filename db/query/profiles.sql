-- name: CreateProfile :one
INSERT INTO profiles (
    id, email, full_name, role, store_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetProfileByEmail :one
SELECT * FROM profiles
WHERE email = $1 LIMIT 1;

-- name: GetProfile :one
SELECT * FROM profiles
WHERE id = $1 LIMIT 1;
