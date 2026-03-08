package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	DB *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS players (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    tier       TEXT NOT NULL CHECK(tier IN ('beginner','explorer','master')),
    avatar     TEXT DEFAULT 'default',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS progress (
    player_id    INTEGER NOT NULL REFERENCES players(id),
    world_id     TEXT NOT NULL,
    mission_id   TEXT NOT NULL,
    stars        INTEGER DEFAULT 0,
    completed_at DATETIME,
    PRIMARY KEY (player_id, world_id, mission_id)
);

CREATE TABLE IF NOT EXISTS world_state (
    player_id  INTEGER NOT NULL REFERENCES players(id),
    world_id   TEXT NOT NULL,
    fs_json    TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (player_id, world_id)
);

CREATE TABLE IF NOT EXISTS history (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id  INTEGER NOT NULL REFERENCES players(id),
    command    TEXT NOT NULL,
    timestamp  DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

func Open(path string) (*DB, error) {
	sqldb, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := sqldb.Exec(schema); err != nil {
		_ = sqldb.Close()
		return nil, err
	}
	return &DB{DB: sqldb}, nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}
