package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/shanewilliams/shell-quest/internal/db"
)

// playerItem implements list.Item for a saved player profile.
type playerItem struct {
	player *db.Player // nil means "New Pirate"
}

func (p playerItem) Title() string {
	if p.player == nil {
		return "[ New Pirate ]"
	}
	return p.player.Name
}

func (p playerItem) Description() string {
	if p.player == nil {
		return "Start a fresh adventure"
	}
	return fmt.Sprintf("Tier: %s", p.player.Tier)
}

func (p playerItem) FilterValue() string { return p.Title() }

func newProfileList(players []*db.Player, w, h int) list.Model {
	items := make([]list.Item, len(players)+1)
	for i, p := range players {
		items[i] = playerItem{player: p}
	}
	items[len(players)] = playerItem{player: nil}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#f5a623")).
		BorderLeftForeground(lipgloss.Color("#f5a623"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#c47d0e")).
		BorderLeftForeground(lipgloss.Color("#f5a623"))

	l := list.New(items, delegate, w, h)
	l.Title = "Choose Your Pirate"
	l.Styles.Title = TitleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	return l
}

func (m Model) profileSelectView() string {
	return m.profileList.View()
}
