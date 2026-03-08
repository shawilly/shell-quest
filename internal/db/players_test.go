package db_test

import (
	"path/filepath"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/db"
)

func testDB(t *testing.T) *db.DB {
	t.Helper()
	d, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

func TestCreatePlayer(t *testing.T) {
	d := testDB(t)
	p, err := d.CreatePlayer("Matilda", "beginner")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Matilda" || p.Tier != "beginner" || p.ID == 0 {
		t.Errorf("unexpected player: %+v", p)
	}
}

func TestListPlayers_EmptyThenOne(t *testing.T) {
	d := testDB(t)
	players, err := d.ListPlayers()
	if err != nil || len(players) != 0 {
		t.Fatalf("expected empty list, got %v %v", players, err)
	}
	if _, err := d.CreatePlayer("Sam", "explorer"); err != nil {
		t.Fatal(err)
	}
	players, err = d.ListPlayers()
	if err != nil || len(players) != 1 {
		t.Fatalf("expected 1 player, got %v %v", players, err)
	}
}
