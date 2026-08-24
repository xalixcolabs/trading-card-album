-- name: CreateContact :one
INSERT OR IGNORE INTO contact (
    user_id, met_user_id, scanned_at
) VALUES (
    ?, ? ,?
) RETURNING *;

-- name: ListContacts :many
SELECT u.id, u.name, u.email, u.github, u.linkedin, u.web, u.description,
       u.is_admin, u.picture, u.created_at, u.updated_at, c.scanned_at
FROM contact c
JOIN user u ON u.id = c.met_user_id
WHERE c.user_id = ?
ORDER BY c.scanned_at DESC;