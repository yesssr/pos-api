-- name: CreateDetailPenjualan :one
INSERT INTO detail_penjualan (
  id_penjualan,
  id_barang,
  jumlah,
  harga_saat_transaksi,
  subtotal
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetDetailPenjualan :many
SELECT
  d.id,
  d.id_penjualan,
  d.id_barang,
  d.jumlah,
  d.harga_saat_transaksi,
  d.subtotal,
  b.nama as nama_barang
FROM detail_penjualan d
JOIN barang b ON d.id_barang = b.id
WHERE d.id_penjualan = $1;
