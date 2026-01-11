-- name: CreateAuthUser :one
INSERT INTO auth.users (
    email, encrypted_password
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetAuthUserByEmail :one
SELECT * FROM auth.users
WHERE email = $1 LIMIT 1;
