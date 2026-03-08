package tui

import (
	"fmt"
	"strings"
)

func (m Model) adventureLogView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("ADVENTURE LOG") + "\n\n")
	if m.player != nil {
		b.WriteString(fmt.Sprintf("Pirate: %s (%s)\n\n", m.player.Name, m.player.Tier))
	}
	if m.runner != nil && m.runner.IsComplete() {
		b.WriteString(SuccessStyle.Render("SKULL ISLAND: COMPLETED") + "\n")
	} else {
		b.WriteString("Skull Island: In Progress\n")
		if m.runner != nil {
			b.WriteString(fmt.Sprintf("  Objective %d of %d\n",
				m.runner.CurrentObjectiveIndex()+1,
				len(m.runner.Mission().Objectives)))
		}
	}
	b.WriteString("\n\nPress ESC or Enter to return to the game.")
	return b.String()
}
