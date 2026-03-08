package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type tierItem struct {
	name       string
	desc       string
	comingSoon bool
}

func (t tierItem) Title() string {
	if t.comingSoon {
		return t.name + "  [ Coming Soon ]"
	}
	return t.name
}

func (t tierItem) Description() string { return t.desc }
func (t tierItem) FilterValue() string { return t.name }

// tierDelegate renders coming-soon items as faint.
type tierDelegate struct {
	list.DefaultDelegate
}

func (d tierDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	t, ok := item.(tierItem)
	if ok && t.comingSoon {
		faint := lipgloss.NewStyle().Faint(true)
		_, _ = fmt.Fprint(w, faint.Render("  "+t.Title()+"\n  "+t.Description()))
		return
	}
	d.DefaultDelegate.Render(w, m, index, item)
}

var allTierItems = []list.Item{
	tierItem{"Beginner", "Ages 3–6 · ls, cd, pwd, cat, echo, clear, help", false},
	tierItem{"Explorer", "Ages 6–8 · + mkdir, touch, cp, mv, rm, find", true},
	tierItem{"Master", "Ages 8–10 · + grep, chmod, man, history, pipes, globs", true},
}

func newTierList(playerName string, w, h int) list.Model {
	delegate := tierDelegate{list.NewDefaultDelegate()}
	delegate.Styles.SelectedTitle =
		delegate.Styles.SelectedTitle.
			Foreground(lipgloss.Color("#f5a623")).
			BorderLeftForeground(lipgloss.Color("#f5a623"))

	l := list.New(allTierItems, delegate, w, h)
	l.Title = "Choose Your Difficulty, " + playerName
	l.Styles.Title = TitleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	return l
}

func (m Model) tierSelectView() string {
	return m.tierList.View()
}
