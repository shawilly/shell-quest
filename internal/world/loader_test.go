package world_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/world"
)

func TestLoadWorld_ParsesMissions(t *testing.T) {
	w, err := world.LoadWorld("skull_island")
	if err != nil {
		t.Fatalf("LoadWorld: %v", err)
	}
	if w.ID != "skull_island" {
		t.Errorf("expected skull_island, got %q", w.ID)
	}
	if len(w.Missions) < 2 {
		t.Errorf("expected at least 2 missions, got %d", len(w.Missions))
	}
}

func TestLoadWorld_MissionHasObjectives(t *testing.T) {
	w, _ := world.LoadWorld("skull_island")
	m := w.Missions[0]
	if len(m.Objectives) == 0 {
		t.Error("expected objectives in mission 1")
	}
	if m.Objectives[0].Command != "ls" {
		t.Errorf("expected first objective to be 'ls', got %q", m.Objectives[0].Command)
	}
}

func TestLoadWorld_MissionHasFilesystem(t *testing.T) {
	w, _ := world.LoadWorld("skull_island")
	m := w.Missions[0]
	if len(m.Filesystem) == 0 {
		t.Error("expected filesystem entries in mission 1")
	}
	entry, ok := m.Filesystem["/docks"]
	if !ok || entry.Type != "dir" {
		t.Errorf("expected /docks to be a dir: %+v", entry)
	}
}

func TestLoadWorld_MissionNarrativeFields(t *testing.T) {
	w, _ := world.LoadWorld("skull_island")
	m := w.Missions[0]
	if m.IntroDialogue == nil {
		t.Fatal("expected intro_dialogue on mission 1")
	}
	if m.IntroDialogue.NPC != "Old Pete" {
		t.Errorf("expected NPC 'Old Pete', got %q", m.IntroDialogue.NPC)
	}
	if m.StartingCWD != "/docks" {
		t.Errorf("expected starting_cwd '/docks', got %q", m.StartingCWD)
	}
	if len(m.ObjectiveHints) == 0 {
		t.Error("expected objective_hints on mission 1")
	}
}

func TestMissionRunner_CurrentHint(t *testing.T) {
	w, _ := world.LoadWorld("skull_island")
	runner := world.NewMissionRunner(w.Missions[0])
	// No hint before any objective completed
	if runner.CurrentHint() != "" {
		t.Errorf("expected empty hint before first objective, got %q", runner.CurrentHint())
	}
}

func TestLoadWorld_UnknownWorld_Errors(t *testing.T) {
	_, err := world.LoadWorld("doesnotexist")
	if err == nil {
		t.Error("expected error for unknown world")
	}
}
