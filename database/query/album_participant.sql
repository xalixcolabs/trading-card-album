-- name: GetAlbumParticipant :one
SELECT * FROM album_participant
WHERE album_id = ? AND user_id = ? LIMIT 1;

-- name: GetCardByAlbumParticipant :one
SELECT c.*
FROM album_participant ap
JOIN card c ON ap.assigned_card_id = c.id
WHERE ap.album_id = ?
  AND ap.assigned_card_id = ?
  AND ap.user_id = ?;

-- name: ListAlbumParticipants :many
SELECT * FROM album_participant
ORDER BY created_at;

-- name: CreateAlbumParticipant :one
INSERT INTO album_participant (
  album_id, user_id, assigned_card_id, joined_at, secret
) VALUES (
  ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateAlbumParticipantSecret :one
UPDATE album_participant
SET secret = ?
WHERE album_id = ? AND user_id = ?
RETURNING *;
