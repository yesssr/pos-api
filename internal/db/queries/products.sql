-- name: CreateProduct :one
INSERT INTO products (name, price, stock, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET
  name = $2,
  price = $3,
  stock = $4,
  image_url = $5,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;

-- name: ListProductAsc :many
SELECT
  id,
  name,
  price,
  stock,
  image_url,
  is_active
FROM products
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: ListProductDesc :many
SELECT
  id,
  name,
  price,
  stock,
  image_url,
  is_active
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
