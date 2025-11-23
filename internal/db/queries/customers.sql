-- name: CreateCustomer :one
INSERT INTO customers (name, phone, address)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers SET
  name = $2,
  phone = $3,
  address = $4
WHERE id = $1
RETURNING *;

-- name: ListCustomerAsc :many
SELECT
  id,
  name,
  phone,
  address,
  created_at
FROM customers
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: ListCustomerDesc :many
SELECT
  id,
  name,
  phone,
  address,
  created_at
FROM customers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetCustomerByID :one
SELECT
  id,
  name,
  phone,
  address,
  created_at
FROM customers
WHERE id = $1;

-- name: DeleteCustomer :one
DELETE FROM customers
WHERE id = $1
RETURNING *;
