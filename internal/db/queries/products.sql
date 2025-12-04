-- name: CreateProduct :one
INSERT INTO products (name, price, stock, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProduct :one
SELECT
  id,
  name,
  price,
  stock,
  image_url,
  is_active
FROM products
WHERE id = $1;

-- name: UpdateProduct :one
UPDATE products SET
  name = $2,
  price = $3,
  stock = $4,
  image_url = $5,
  is_active = $6,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;

-- name: ListProductsAsc :many
SELECT
  id,
  name,
  price,
  stock,
  image_url,
  is_active
FROM products
WHERE name ILIKE '%' || $4 || '%'
ORDER BY $3 ASC
LIMIT $1 OFFSET $2;

-- name: ListProductsDesc :many
SELECT
  id,
  name,
  price,
  stock,
  image_url,
  is_active
FROM products
WHERE name ILIKE '%' || $4 || '%'
ORDER BY $3 DESC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) AS count
FROM products
WHERE name ILIKE '%' || $1 || '%';

-- name: ListProductsActive :many
SELECT
  id,
  name,
  price,
  stock,
  image_url,
  is_active
FROM products
WHERE is_active = TRUE
AND name ILIKE '%' || $4 || '%'
ORDER BY $3 DESC
LIMIT $1 OFFSET $2;

-- name: CountProductsActive :one
SELECT COUNT(*) AS count
FROM products
WHERE is_active = TRUE
AND name ILIKE '%' || $1 || '%';


-- name: UpdateProductStock :one
UPDATE products SET
  stock = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING id;

-- name: GetProductForUpdate :one
SELECT id, name, stock, price
FROM products
WHERE id = $1
FOR UPDATE;
