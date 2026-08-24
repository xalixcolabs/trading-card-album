-- migrate:up
CREATE TABLE IF NOT EXISTS card(
    id TEXT PRIMARY KEY,
    album_id TEXT NOT NULL,
    number TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    FOREIGN KEY (album_id) REFERENCES albums(album_id) ON DELETE CASCADE,
    UNIQUE (album_id, number)
);

-- migrate:down
DROP TABLE IF NOT EXISTS card;
