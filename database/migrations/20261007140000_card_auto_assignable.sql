-- migrate:up
ALTER TABLE card ADD COLUMN auto_assignable INTEGER NOT NULL DEFAULT 1;

-- migrate:down
ALTER TABLE card DROP COLUMN auto_assignable;
