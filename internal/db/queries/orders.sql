-- name: GetUserOrders :many
SELECT * FROM orders
WHERE user_id=$1
ORDER BY created_at DESC;

-- name: CreateOrder :one
INSERT INTO orders(
    user_id, total_price
) VALUES (
    $1, $2
)
RETURNING *;
