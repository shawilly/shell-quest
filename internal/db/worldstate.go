package db

import "database/sql"

// SaveWorldState persists the FS JSON for a player+world combo (upsert).
func (d *DB) SaveWorldState(playerID int64, worldID, fsJSON string) error {
	_, err := d.DB.Exec(`
		INSERT INTO world_state (player_id, world_id, fs_json, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(player_id, world_id) DO UPDATE
		SET fs_json = excluded.fs_json, updated_at = excluded.updated_at
	`, playerID, worldID, fsJSON)
	return err
}

// LoadWorldState retrieves the saved FS JSON for a player+world combo.
// Returns sql.ErrNoRows if no state exists yet.
func (d *DB) LoadWorldState(playerID int64, worldID string) (string, error) {
	row := d.DB.QueryRow(
		`SELECT fs_json FROM world_state WHERE player_id = ? AND world_id = ?`,
		playerID, worldID,
	)
	var fsJSON string
	if err := row.Scan(&fsJSON); err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}
	return fsJSON, nil
}
