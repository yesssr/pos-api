-- name: CreateDetailTransaction :one
INSERT INTO detail_transactions (
  id_transaction,
  id_product,
  qty,
  price,
  subtotal
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetDetailTransaction :many
SELECT
  d.id,
  d.id_transaction,
  d.id_product,
  d.qty,
  d.price,
  d.subtotal,
  p.name as product_name
FROM detail_transactions d
JOIN products p ON d.id_product = p.id
WHERE d.id_transaction = $1;
