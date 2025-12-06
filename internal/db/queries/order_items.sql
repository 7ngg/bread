-- name: GetOrderItems :many
SELECT * FROM order_items
WHERE order_id=$1;

-- name: InsertOrderItems :copyfrom
INSERT INTO order_items(
    order_id, product_id, quantity
) VALUES (
    $1, $2, $3
);

