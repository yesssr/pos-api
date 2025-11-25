-- name: CreateTransaction :one
INSERT INTO transactions (
  id_user,
  id_customer,
  total,
  payment_method,
  payment_status,
  id_transaction_gateway
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListTransactions :many
SELECT
  t.id,
  t.id_user,
  t.id_customer,
  t.total,
  t.date,
  t.payment_method,
  t.payment_status,
  t.id_transaction_gateway,
  c.name AS customer_name,
  u.username AS cashier
FROM transactions t
LEFT JOIN customers c ON t.id_customer = c.id
JOIN users u ON t.id_user = u.id
WHERE date BETWEEN $1 AND $2
ORDER BY t.date DESC
LIMIT $3 OFFSET $4;

-- name: CountTransactions :one
SELECT COUNT(*) AS count
FROM transactions
WHERE is_active = true;
