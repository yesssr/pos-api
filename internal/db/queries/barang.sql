-- name: CreateBarang :one
INSERT INTO barang (nama, harga, stok, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateBarang :one
UPDATE barang SET
  nama = $2,
  harga = $3,
  stok = $4,
  image_url = $5,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBarang :one
DELETE FROM barang
WHERE id = $1
RETURNING *;

-- name: ListBarangAsc :many
SELECT
  id,
  nama,
  harga,
  stok,
  image_url,
  is_active
FROM barang
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: ListBarangDesc :many
SELECT
  id,
  nama,
  harga,
  stok,
  image_url,
  is_active
FROM barang
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
