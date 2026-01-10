-- name: CreateCategory :one
INSERT INTO categories (outlet_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories
WHERE outlet_id = $1
ORDER BY name;

-- name: CreateProduct :one
INSERT INTO products (
    outlet_id, category_id, name, price, stock, image_url, is_available
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE outlet_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateProductStock :exec
UPDATE products
SET stock = stock + $2, updated_at = NOW()
WHERE id = $1;
