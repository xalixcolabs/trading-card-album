-- migrate:up
CREATE TABLE IF NOT EXISTS user(
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    github TEXT NOT NULL,
    linkedin TEXT NOT NULL,
    web TEXT NOT NULL,
    description TEXT NOT NULL,
    is_admin INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS user;
