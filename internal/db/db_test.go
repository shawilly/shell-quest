package db_test

import (
	"path/filepath"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/db"
)

func TestOpen_CreatesSchema(t *testing.T) {
	dir := t.TempDir()
	d, err := db.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	// Verify tables exist by querying them
	tables := []string{"players", "progress", "world_state", "history"}
	for _, tbl := range tables {
		row := d.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tbl)
		var name string
		if err := row.Scan(&name); err != nil {
			t.Errorf("table %q not found: %v", tbl, err)
		}
	}
}
