-- name: CreatePelanggan :one
INSERT INTO pelanggan (nama, no_hp, alamat)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdatePelanggan :one
UPDATE pelanggan SET
  nama = $2,
  no_hp = $3,
  alamat = $4
WHERE id = $1
RETURNING *;

-- name: ListPelangganAsc :many
SELECT
  id,
  nama,
  no_hp,
  alamat,
  created_at
FROM pelanggan
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: ListPelangganDesc :many
SELECT
  id,
  nama,
  no_hp,
  alamat,
  created_at
FROM pelanggan
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPelangganByID :one
SELECT
  id,
  nama,
  no_hp,
  alamat,
  created_at
FROM pelanggan
WHERE id = $1;

-- name: DeletePelanggan :one
DELETE FROM pelanggan
WHERE id = $1
RETURNING *;
