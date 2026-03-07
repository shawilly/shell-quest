package db_test

import (
	"testing"
)

func TestSaveAndLoadWorldState(t *testing.T) {
	d := testDB(t) // reuses the helper from players_test.go
	player, _ := d.CreatePlayer("TestPirate", "beginner")

	fsJSON := `{"test": "data"}`
	if err := d.SaveWorldState(player.ID, "skull_island", fsJSON); err != nil {
		t.Fatalf("SaveWorldState: %v", err)
	}

	loaded, err := d.LoadWorldState(player.ID, "skull_island")
	if err != nil {
		t.Fatalf("LoadWorldState: %v", err)
	}
	if loaded != fsJSON {
		t.Errorf("expected %q, got %q", fsJSON, loaded)
	}
}

func TestSaveWorldState_Upsert(t *testing.T) {
	d := testDB(t)
	player, _ := d.CreatePlayer("TestPirate2", "beginner")

	d.SaveWorldState(player.ID, "skull_island", `{"v": "1"}`)
	d.SaveWorldState(player.ID, "skull_island", `{"v": "2"}`) // upsert

	loaded, err := d.LoadWorldState(player.ID, "skull_island")
	if err != nil {
		t.Fatal(err)
	}
	if loaded != `{"v": "2"}` {
		t.Errorf("expected v=2, got %q", loaded)
	}
}

func TestLoadWorldState_Missing_Errors(t *testing.T) {
	d := testDB(t)
	player, _ := d.CreatePlayer("TestPirate3", "beginner")
	_, err := d.LoadWorldState(player.ID, "nonexistent_world")
	if err == nil {
		t.Error("expected error for missing world state")
	}
}
