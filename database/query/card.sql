-- name: GetCard :one
SELECT * FROM card
WHERE id = ? LIMIT 1;

-- name: ListCards :many
SELECT * FROM card
ORDER BY number;

-- name: ListCardsByAlbumId :many
SELECT * FROM card
WHERE album_id = ?
ORDER BY number;

-- name: CreateCard :one
INSERT INTO card (
  id, album_id, number, name, description, image_url, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateCard :one
UPDATE card
SET album_id = ?,
    number = ?,
    name = ?,
    description = ?,
    image_url = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteCard :exec
DELETE FROM card
WHERE id = ?;
-- name: DeleteCardsByAlbumId :exec
DELETE FROM card
WHERE album_id = ?;
