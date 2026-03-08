package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func newShellInput() textinput.Model {
	ti := textinput.New()
	ti.Prompt = "" // prompt is prepended manually in shellContent so it can include cwd
	ti.Width = 60
	ti.TextStyle = OutputStyle
	return ti
}

func newShellViewport(w, h int) viewport.Model {
	vp := viewport.New(w, h)
	vp.SetContent("")
	return vp
}

func (m Model) gameView() string {
	if m.width == 0 {
		return "Loading Shell Quest..."
	}

	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth
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
	if m.runner != nil && !m.runner.IsComplete() {
		if obj := m.runner.CurrentObjective(); obj != nil {
			b.WriteString(fmt.Sprintf("\n\nObjective %d: use '%s'",
				m.runner.CurrentObjectiveIndex()+1, obj.Command))
		}
	}
	return b.String()
}

func (m Model) shellContent() string {
	prompt := PromptStyle.Render("pirate@quest:" + m.cwd + "$ ")
	input := prompt + m.shellInput.View() + "_"

	return TitleStyle.Render("SHELL QUEST") + "\n\n" +
		m.shellVP.View() + "\n" +
		input
}
