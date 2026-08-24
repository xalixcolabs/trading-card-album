-- name: GetUser :one
SELECT * FROM user
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM user
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM user
ORDER BY name;

-- name: SearchUsersByEmail :many
SELECT * FROM user
WHERE email LIKE ?
ORDER BY name;

-- name: CreateUser :one
INSERT INTO user (
  id, name, email, github, linkedin, web, description, is_admin, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) ON CONFLICT(email) DO UPDATE SET
    name = excluded.name
RETURNING *;

-- name: UpdateUser :one
UPDATE user
SET name = ?,
email = ?,
github = ?,
linkedin = ?,
web = ?,
description = ?,
updated_at = ?
WHERE id = ?
RETURNING *;

-- name: UpdateUserIsAdmin :one
UPDATE user
SET is_admin = ?,
updated_at = ?
WHERE id = ?
RETURNING *;
