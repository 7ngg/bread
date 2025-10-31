-- name: GetProducts :many
SELECT * FROM products
ORDER BY name
LIMIT ? OFFSET ?;

-- name: ProductsCount :one
SELECT COUNT(id)
FROM products;
