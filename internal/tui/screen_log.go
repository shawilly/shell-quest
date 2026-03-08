package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	listy "github.com/genekkion/theHermit/list"
)

// hermitItem implements listy.Item for Hermit overlay lines.
type hermitItem struct {
	title string
}

func (h hermitItem) Title() string { return h.title }

func newHermit(height, width int) listy.Model {
	l := listy.New(height, width, nil)
	l.SetTitle("ADVENTURE LOG")
	l.SetBorderForeground(lipgloss.Color("#f5a623"))
	l.SetIsNumbered(false)
	return l
}

func adventureLogItems(m Model) []listy.Item {
	var lines []string
	if m.player != nil {
		lines = append(lines, fmt.Sprintf("Pirate: %s (%s)", m.player.Name, m.player.Tier))
		lines = append(lines, "")
	}
	if m.runner != nil && m.runner.IsComplete() {
		lines = append(lines, "✓ Skull Island: COMPLETED")
	} else {
		lines = append(lines, "Skull Island: In Progress")
		if m.runner != nil {
			lines = append(lines, fmt.Sprintf("  Objective %d of %d",
				m.runner.CurrentObjectiveIndex()+1,
				len(m.runner.Mission().Objectives)))
		}
	}
	lines = append(lines, "", "ESC / Ctrl+L to close")

	items := make([]listy.Item, len(lines))
	for i, l := range lines {
		items[i] = hermitItem{title: l}
	}
	return items
}

func parentModeItems(m Model) []listy.Item {
	var lines []string
	if !m.parentUnlocked {
		lines = []string{
			fmt.Sprintf("What is %d + %d?", m.mathA, m.mathB),
			"",
			m.mathInput.View(),
			"",
			"ESC to cancel",
		}
	} else {
		lines = []string{
			"Access granted!",
			"",
			"Q - Quit game",
			"ESC - Return to game",
		}
	}
	items := make([]listy.Item, len(lines))
	for i, l := range lines {
		items[i] = hermitItem{title: l}
	}
	return items
}
