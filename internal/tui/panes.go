package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

func (m Model) tierSelectView() string {
	tiers := []struct {
		name       string
		desc       string
		comingSoon bool
	}{
		{"Beginner", "Ages 3-6. Commands: ls, cd, pwd, cat, echo, clear, help", false},
		{"Explorer", "Ages 6-8. + mkdir, touch, cp, mv, rm, find", true},
		{"Master", "Ages 8-10. + grep, chmod, man, history, pipes, globs", true},
	}
	var b strings.Builder
	b.WriteString(TitleStyle.Render("Choose Your Difficulty, "+m.nameInput.Value()) + "\n\n")
	for i, t := range tiers {
		cursor := "  "
		if i == m.selectedIdx {
			cursor = "> "
		}
		if t.comingSoon {
			line := fmt.Sprintf("%s%s  [ Coming Soon ]\n    %s\n\n", cursor, t.name, t.desc)
			b.WriteString(lipgloss.NewStyle().Faint(true).Render(line))
		} else {
			b.WriteString(fmt.Sprintf("%s%s\n    %s\n\n", cursor, t.name, t.desc))
		}
	}
	b.WriteString("Use up/down arrows and Enter to select.")
	return b.String()
}
