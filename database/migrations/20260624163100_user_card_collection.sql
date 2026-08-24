-- migrate:up
CREATE TABLE IF NOT EXISTS user_card_collection (
    user_id TEXT NOT NULL,
    album_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    unlocked_at INTEGER NOT NULL,
    PRIMARY KEY (user_id, album_id, card_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (album_id) REFERENCES album(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_collection_progress 
ON user_card_collection(user_id, album_id);

-- migrate:down
DROP TABLE IF EXISTS user_card_collection;
