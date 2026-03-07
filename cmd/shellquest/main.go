package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shanewilliams/shell-quest/internal/db"
	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/tui"
	"github.com/shanewilliams/shell-quest/internal/world"
)

func main() {
	// DB setup
	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".shellquest")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal(err)
	}
	database, err := db.Open(filepath.Join(dbDir, "progress.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Load world
	w, err := world.LoadWorld("skull_island")
	if err != nil {
		log.Fatal(err)
	}

	// Use first mission for now (profile/tier selection in Task 30)
	mission := w.Missions[0]

	// Build virtual FS from mission filesystem definition
	// Sort paths so parents are created before children
	fs := shell.NewFS()
	paths := make([]string, 0, len(mission.Filesystem))
	for p := range mission.Filesystem {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		entry := mission.Filesystem[p]
		if entry.Type == "dir" {
			fs.Mkdir(p, entry.Hidden)
		} else {
			fs.WriteFile(p, entry.Content, entry.Hidden)
		}
	}

	// Default player for now
	player := &db.Player{Name: "Pirate", Tier: "beginner"}

	// Executor
	ex := shell.NewExecutor(fs)
	tui.RegisterCommands(ex, player.Tier)

	// Mission runner
	runner := world.NewMissionRunner(mission)

	// Start game
	model := tui.NewGameModel(database, player, fs, ex, runner, "/island", mission.StartingClue)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
