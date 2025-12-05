-- name: CreateTransaction :one
INSERT INTO transactions (
  id_user,
  id_customer,
  total,
  payment_method,
  id_transaction_gateway,
  date
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
ORDER BY t.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountTransactions :one
SELECT COUNT(*) AS count
FROM transactions;

-- name: UpdateTransactionStatus :one
UPDATE transactions SET
  payment_status = $2,
  payment_method = $3,
  id_transaction_gateway = $4,
  total = $5,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateStatusByPaymentId :one
UPDATE transactions SET
  payment_status = $2,
  updated_at = NOW()
WHERE id_transaction_gateway = $1
RETURNING *;

-- name: ListSalesPerDay :many
SELECT
  TO_CHAR(date_trunc('day', date), 'YYYY-MM-DD') AS name,
  SUM(total) AS sales
FROM transactions
WHERE date >= $1 AND date <= $2
GROUP BY date_trunc('day', date)
ORDER BY name;

-- name: ListSalesPerWeek :many
SELECT
  TO_CHAR(date_trunc('week', date), 'YYYY-MM-DD') AS name,
  SUM(total) AS sales
FROM transactions
WHERE date >= $1 AND date <= $2
GROUP BY date_trunc('week', date)
ORDER BY name;

-- name: ListSalesPerMonth :many
SELECT
  TO_CHAR(date_trunc('month', date), 'YYYY-MM') AS name,
  SUM(total) AS sales
FROM transactions
WHERE date >= $1 AND date <= $2
GROUP BY date_trunc('month', date)
ORDER BY name;

-- name: ListSalesPerYear :many
SELECT
  TO_CHAR(date, 'YYYY') AS name,
  SUM(total) AS sales
FROM transactions
WHERE date >= $1 AND date <= $2
GROUP BY TO_CHAR(date, 'YYYY')
ORDER BY name;
