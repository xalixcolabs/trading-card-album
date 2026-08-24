-- migrate:up
CREATE TABLE IF NOT EXISTS card_pool (
    album_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    is_drawn INTEGER NOT NULL DEFAULT 0, -- 0 = Disponible en la baraja, 1 = Ya repartida
    FOREIGN KEY (album_id) REFERENCES album(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

CREATE INDEX idx_pool_available ON card_pool(album_id, is_drawn);

-- migrate:down
DROP TABLE IF EXISTS card_pool
