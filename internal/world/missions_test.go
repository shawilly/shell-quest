package world_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/world"
)

func testMission() world.Mission {
	return world.Mission{
		ID:    "m1",
		Title: "Test Mission",
		Objectives: []world.Objective{
			{Command: "ls", Path: "/island"},
			{Command: "cd", Path: "/island/cave"},
		},
	}
}

func TestMissionRunner_InitiallyNotComplete(t *testing.T) {
	r := world.NewMissionRunner(testMission())
	if r.IsComplete() {
		t.Error("should not be complete at start")
	}
	if r.CurrentObjectiveIndex() != 0 {
		t.Error("should start at objective 0")
	}
}

func TestMissionRunner_AdvancesOnMatchingEvent(t *testing.T) {
	r := world.NewMissionRunner(testMission())
	advanced := r.HandleEvent(&shell.Event{Type: "ls", Path: "/island"})
	if !advanced {
		t.Error("expected HandleEvent to return true on match")
	}
	if r.CurrentObjectiveIndex() != 1 {
		t.Error("should be at objective 1")
	}
}

func TestMissionRunner_CompletesAfterAllObjectives(t *testing.T) {
	r := world.NewMissionRunner(testMission())
	r.HandleEvent(&shell.Event{Type: "ls", Path: "/island"})
	r.HandleEvent(&shell.Event{Type: "cd", Path: "/island/cave"})
	if !r.IsComplete() {
		t.Error("should be complete after all objectives met")
	}
}

func TestMissionRunner_WrongEvent_DoesNotAdvance(t *testing.T) {
	r := world.NewMissionRunner(testMission())
	advanced := r.HandleEvent(&shell.Event{Type: "cat", Path: "/island/note.txt"})
	if advanced {
		t.Error("wrong event should not advance")
	}
	if r.CurrentObjectiveIndex() != 0 {
		t.Error("should still be at objective 0")
	}
}

func TestMissionRunner_NilEvent_DoesNotAdvance(t *testing.T) {
	r := world.NewMissionRunner(testMission())
	advanced := r.HandleEvent(nil)
	if advanced {
		t.Error("nil event should not advance")
	}
}

func TestMissionRunner_CurrentObjective_NilWhenComplete(t *testing.T) {
	r := world.NewMissionRunner(testMission())
	r.HandleEvent(&shell.Event{Type: "ls", Path: "/island"})
	r.HandleEvent(&shell.Event{Type: "cd", Path: "/island/cave"})
	if r.CurrentObjective() != nil {
		t.Error("CurrentObjective should be nil when complete")
	}
}
