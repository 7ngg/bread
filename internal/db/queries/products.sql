-- name: GetProducts :many
SELECT * FROM products
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ProductsCount :one
SELECT COUNT(id)
FROM products;
