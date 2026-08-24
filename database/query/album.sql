-- name: GetAlbum :one
SELECT * FROM album
WHERE id = ? LIMIT 1;

-- name: ListAlbums :many
SELECT * FROM album
ORDER BY created_at;

-- name: ListAlbumsByUserId :many
SELECT a.*
FROM album_participant ap
JOIN album a ON ap.album_id = a.id
WHERE ap.user_id = ?
ORDER BY ap.joined_at DESC;

-- name: CreateAlbum :one
INSERT INTO album (
  id, title, total_cards, created_at
) VALUES (
  ?, ?, ?, ?
) RETURNING *;

-- name: UpdateAlbum :one
UPDATE album
SET title = ?,
total_cards = ?
WHERE id = ?
RETURNING *;
