-- migrate:up
CREATE TABLE IF NOT EXISTS contact(
    user_id TEXT NOT NULL,
    met_user_id TEXT NOT NULL,
    scanned_at INTEGER NOT NULL,
    PRIMARY KEY (user_id, met_user_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (met_user_id) REFERENCES user(id) ON DELETE CASCADE,
    CHECK (user_id <> met_user_id)
);

CREATE INDEX idx_user_contacts_owner ON contact(user_id);

-- migrate:down
DROP TABLE IF EXISTS contact;

