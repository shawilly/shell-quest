package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shanewilliams/shell-quest/internal/db"
	"github.com/shanewilliams/shell-quest/internal/tui"
)

func main() {
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

	model := tui.NewStartupModel(database)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
