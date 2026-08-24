-- name: CollectCard :exec
INSERT OR IGNORE INTO user_card_collection (user_id, album_id, card_id, unlocked_at)
VALUES (?, ?, ?, ?);

-- name: GetUserCollection :many
SELECT c.*
FROM user_card_collection ucc
JOIN card c ON ucc.card_id = c.id
WHERE ucc.user_id = ? AND ucc.album_id = ?
ORDER BY c.number;

-- name: ListCollectedCardsByUser :many
SELECT c.*
FROM user_card_collection ucc
JOIN card c ON ucc.card_id = c.id
WHERE ucc.user_id = ?
ORDER BY c.album_id, c.number;

-- name: CardInCollection :one
SELECT EXISTS(
    SELECT 1 FROM user_card_collection
    WHERE user_id = ? AND album_id = ? AND card_id = ?
) AS owned;

-- name: DeleteUserCollectionByAlbumId :exec
DELETE FROM user_card_collection
WHERE album_id = ?;
