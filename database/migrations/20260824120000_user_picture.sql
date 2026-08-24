-- migrate:up
ALTER TABLE user ADD COLUMN picture TEXT NOT NULL DEFAULT '';

-- migrate:down
ALTER TABLE user DROP COLUMN picture;