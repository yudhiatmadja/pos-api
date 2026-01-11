-- name: CreateOrder :one
INSERT INTO orders (
    store_id, table_session_id, cashier_id, order_number, 
    total_amount, tax_amount, discount_amount, final_amount, note, status, payment_status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, product_id, product_name, product_price, quantity, total_price, note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrdersByStore :many
SELECT * FROM orders
WHERE store_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrderPaymentStatus :one
UPDATE orders
SET payment_status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetOrdersBySession :many
SELECT * FROM orders
WHERE table_session_id = $1
ORDER BY created_at DESC;
