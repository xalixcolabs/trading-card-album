-- name: GetRandomAvailableCard :one
SELECT card_id 
FROM card_pool 
WHERE album_id = ? AND is_drawn = 0 
ORDER BY RANDOM() 
LIMIT 1;

-- name: ResetCardPool :one
UPDATE card_pool 
SET is_drawn = 0 
WHERE album_id = ?
RETURNING *;

-- name: MarkCardAsDrawn :one
UPDATE card_pool 
SET is_drawn = 1 
WHERE album_id = ? AND card_id = ?
RETURNING *;

-- name: CreateCardPoolRow :one
INSERT INTO card_pool(
    album_id, card_id
) VALUES (
    ?, ?
) RETURNING *;
-- name: DeleteCardPoolByAlbumId :exec
DELETE FROM card_pool
WHERE album_id = ?;
