package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#f5a623"))
	return s
}

func (m Model) loadingView() string {
	return "\n\n  " + m.spinner.View() + "  Loading your adventure...\n\n" +
		"  (Skull Island is being charted)"
}
