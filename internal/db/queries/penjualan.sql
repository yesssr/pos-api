-- name: CreatePenjualan :one
INSERT INTO penjualan (
  id_user,
  id_pelanggan,
  total,
  metode_pembayaran,
  status_pembayaran,
  id_transaksi_gateway
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListPenjualan :many
SELECT
  p.id,
  p.id_user,
  p.id_pelanggan,
  p.total,
  p.tanggal,
  p.metode_pembayaran,
  p.status_pembayaran,
  p.id_transaksi_gateway,
  pl.nama AS nama_pelanggan,
  u.username AS kasir
FROM penjualan p
LEFT JOIN pelanggan pl ON p.id_pelanggan = pl.id
JOIN users u ON p.id_user = u.id
WHERE tanggal BETWEEN $1 AND $2
ORDER BY p.tanggal DESC
LIMIT $3 OFFSET $4;
