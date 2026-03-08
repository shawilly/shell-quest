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

func (m Model) gameView() string {
	if m.width == 0 {
		// Not yet sized — return minimal view
		return "Loading Shell Quest..."
	}

	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	// Inner widths account for border + padding
	storyInner := leftWidth - 8
	shellInner := rightWidth - 8
	if storyInner < 10 {
		storyInner = 10
	}
	if shellInner < 10 {
		shellInner = 10
	}

	storyPane := StoryBorder.Width(storyInner).Render(m.storyContent())
	shellPane := ShellBorder.Width(shellInner).Render(m.shellContent())

	return lipgloss.JoinHorizontal(lipgloss.Top, storyPane, shellPane)
}

func (m Model) storyContent() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("TREASURE MAP") + "\n\n")

	if m.storyText != "" {
		b.WriteString(m.storyText + "\n\n")
	}

	if m.clueText != "" {
		b.WriteString("CLUE:\n" + m.clueText)
	}

	// Objective progress
	if m.runner != nil && !m.runner.IsComplete() {
		obj := m.runner.CurrentObjective()
		if obj != nil {
			b.WriteString(fmt.Sprintf("\n\nObjective %d: use '%s'",
				m.runner.CurrentObjectiveIndex()+1, obj.Command))
		}
	}

	return b.String()
}

func (m Model) shellContent() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("SHELL QUEST") + "\n\n")

	for _, line := range m.outputLines {
		b.WriteString(line + "\n")
	}

	// Prompt
	prompt := PromptStyle.Render("pirate@quest:" + m.cwd + "$ ")
	b.WriteString(prompt + m.inputBuf + "_")

	return b.String()
}

func (m Model) profileSelectView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("SHELL QUEST - Choose Your Pirate") + "\n\n")

	for i, p := range m.profiles {
		cursor := "  "
		if i == m.selectedIdx {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s (%s)\n", cursor, p.Name, p.Tier))
	}

	// New profile option
	cursor := "  "
	if m.selectedIdx == len(m.profiles) {
		cursor = "> "
	}
	b.WriteString(cursor + "[ New Pirate ]\n")
	b.WriteString("\nUse up/down arrows and Enter to select.")
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
