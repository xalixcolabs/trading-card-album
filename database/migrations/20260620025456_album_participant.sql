-- migrate:up
CREATE TABLE IF NOT EXISTS album_participant(
    album_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    assigned_card_id TEXT NOT NULL,
    joined_at INTEGER NOT NULL,
    secret TEXT NOT NULL,
    PRIMARY KEY (album_id, user_id),
    FOREIGN KEY (album_id) REFERENCES album(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_card_id) REFERENCES card(id)
);

-- migrate:down
DROP TABLE IF EXISTS album_participant;
