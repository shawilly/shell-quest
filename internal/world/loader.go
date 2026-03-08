package world

import (
	"encoding/json"
	"fmt"

	"github.com/shanewilliams/shell-quest/content"
)

var worldFiles = map[string]string{
	"skull_island": "worlds/skull_island.json",
}

// FSEntry represents a file or directory in a mission's filesystem definition.
type FSEntry struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Hidden  bool   `json:"hidden"`
}

// Objective is a goal within a mission (a command + path that must be executed).
type Objective struct {
	Command string `json:"command"`
	Path    string `json:"path"`
}

// NPC holds an NPC name and their dialogue text.
type NPC struct {
	NPC  string `json:"npc"`
	Text string `json:"text"`
}

// Mission represents a single treasure hunt mission.
type Mission struct {
	ID             string              `json:"id"`
	Title          string              `json:"title"`
	StartingClue   string              `json:"starting_clue"`
	StartingCWD    string              `json:"starting_cwd,omitempty"`
	Objectives     []Objective         `json:"objectives"`
	SuccessMessage string              `json:"success_message"`
	Treasure       string              `json:"treasure"`
	Unlocks        []string            `json:"unlocks"`
	Filesystem     map[string]*FSEntry `json:"filesystem"`
	IntroDialogue  *NPC                `json:"intro_dialogue,omitempty"`
	ObjectiveHints []string            `json:"objective_hints,omitempty"`
	BugTaunt       string              `json:"bug_taunt,omitempty"`
}

// World represents a game world containing multiple missions.
type World struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Missions []Mission `json:"missions"`
}

// LoadWorld loads a world by ID from embedded content.
func LoadWorld(id string) (*World, error) {
	path, ok := worldFiles[id]
	if !ok {
		return nil, fmt.Errorf("world not found: %s", id)
	}
	data, err := content.WorldsFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read world %s: %w", id, err)
	}
	var w World
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("parse world %s: %w", id, err)
	}
	return &w, nil
}
