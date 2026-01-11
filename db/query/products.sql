-- name: CreateProduct :one
INSERT INTO products (
    store_id, category_id, name, description, sku, price, stock, image_url, is_available
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE store_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateProductStock :one
UPDATE products
SET stock = stock + $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
