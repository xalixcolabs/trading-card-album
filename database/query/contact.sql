-- name: CreateContact :one
INSERT OR IGNORE INTO contact (
    user_id, met_user_id, scanned_at
) VALUES (
    ?, ? ,?
) RETURNING *;
