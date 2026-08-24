CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE user(
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    github TEXT NOT NULL,
    linkedin TEXT NOT NULL,
    web TEXT NOT NULL,
    description TEXT NOT NULL,
    is_admin INTEGER NOT NULL DEFAULT 0,
    picture TEXT NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
CREATE TABLE album(
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    total_cards INTEGER NOT NULL,
    created_at INTEGER NOT NULL
);
CREATE TABLE card(
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
CREATE TABLE album_participant(
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
CREATE TABLE contact(
    user_id TEXT NOT NULL,
    met_user_id TEXT NOT NULL,
    scanned_at INTEGER NOT NULL,
    PRIMARY KEY (user_id, met_user_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (met_user_id) REFERENCES user(id) ON DELETE CASCADE,
    CHECK (user_id <> met_user_id)
);
CREATE INDEX idx_user_contacts_owner ON contact(user_id);
CREATE TABLE card_pool (
    album_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    is_drawn INTEGER NOT NULL DEFAULT 0, -- 0 = Disponible en la baraja, 1 = Ya repartida
    FOREIGN KEY (album_id) REFERENCES album(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);
CREATE INDEX idx_pool_available ON card_pool(album_id, is_drawn);
CREATE TABLE user_card_collection (
    user_id TEXT NOT NULL,
    album_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    unlocked_at INTEGER NOT NULL,
    PRIMARY KEY (user_id, album_id, card_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (album_id) REFERENCES album(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_collection_progress
ON user_card_collection(user_id, album_id);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20260620020715'),
  ('20260620021011'),
  ('20260620024528'),
  ('20260620025456'),
  ('20260620030252'),
  ('20260624144904'),
  ('20260624163100'),
  ('20260824120000');
