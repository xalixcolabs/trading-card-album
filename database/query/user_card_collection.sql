-- name: CollectCard :exec
INSERT OR IGNORE INTO user_card_collection (user_id, album_id, card_id, unlocked_at)
VALUES (?, ?, ?, ?);

-- name: GetUserCollection :many
SELECT c.*
FROM user_card_collection ucc
JOIN card c ON ucc.card_id = c.id
WHERE ucc.user_id = ? AND ucc.album_id = ?
ORDER BY c.number;
