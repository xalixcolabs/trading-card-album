-- migrate:up
ALTER TABLE user ADD COLUMN public_email INTEGER NOT NULL DEFAULT 0;
ALTER TABLE user ADD COLUMN public_contact TEXT NOT NULL DEFAULT '';

-- migrate:down
ALTER TABLE user DROP COLUMN public_contact;
ALTER TABLE user DROP COLUMN public_email;
