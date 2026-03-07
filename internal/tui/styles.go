package tui

import "github.com/charmbracelet/lipgloss"

var (
	StoryBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f5a623")).
		Padding(1, 2)

	ShellBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a90d9")).
		Padding(1, 2)

	PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f5a623")).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff4444"))

	OutputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e0e0e0"))

	TitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f5a623")).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4caf50")).
		Bold(true)
)
