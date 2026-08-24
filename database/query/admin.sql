-- name: GetOverviewStats :one
SELECT
    (SELECT COUNT(*) FROM album) AS albums,
    (SELECT COUNT(*) FROM user) AS users,
    (SELECT COUNT(*) FROM card) AS cards,
    (SELECT COUNT(*) FROM album_participant) AS participants,
    (SELECT COUNT(*) FROM contact) AS contacts,
    (SELECT COUNT(*) FROM user_card_collection) AS collected;