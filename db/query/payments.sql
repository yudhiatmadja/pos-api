-- name: CreatePayment :one
INSERT INTO payments (
    order_id, payment_method, amount, status, qris_url
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdatePaymentQRIS :exec
UPDATE payments
SET qris_url = $2
WHERE id = $1;

-- name: GetPaymentByOrder :one
SELECT * FROM payments
WHERE order_id = $1 LIMIT 1;
