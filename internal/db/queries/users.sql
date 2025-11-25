-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  role,
  image_url
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListUsers :many
SELECT
  id,
  username,
  role,
  image_url
FROM users
WHERE is_active = true
LIMIT $1 OFFSET $2;

-- name: GetUserById :one
SELECT
  id,
  username,
  role,
  image_url
FROM users
WHERE id = $1
AND is_active = true;

-- name: UpdateUser :one
UPDATE users SET
  username = $2,
  role = $3,
  is_active = COALESCE($4, is_active),
  image_url = $5,
  updated_at = NOW()
WHERE id = $1
AND is_active = true
RETURNING *;

-- name: UpdatePass :one
UPDATE users SET
  password = $2,
  updated_at = NOW()
WHERE id = $1
AND is_active = true
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;

-- name: GetUserByUsername :one
SELECT
  id,
  username,
  password,
  role,
  image_url
FROM users
WHERE username = $1
AND is_active = true;

-- name: CountUsers :one
SELECT COUNT(*) AS count
FROM users
WHERE is_active = true;
