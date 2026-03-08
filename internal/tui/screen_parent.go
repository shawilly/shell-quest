package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
)

func newMathInput() textinput.Model {
	ti := textinput.New()
	ti.CharLimit = 3
	ti.Width = 5
	ti.Prompt = "> "
	ti.PromptStyle = PromptStyle
	ti.Validate = func(s string) error {
		for _, r := range s {
			if r < '0' || r > '9' {
				return fmt.Errorf("digits only")
			}
		}
		return nil
	}
	return ti
}

func (m Model) parentModeView() string {
	if !m.parentUnlocked {
		return TitleStyle.Render("PARENT MODE") + "\n\n" +
			fmt.Sprintf("What is %d + %d?\n\n", m.mathA, m.mathB) +
			m.mathInput.View() + "\n\n" +
			"Press ESC to cancel."
	}
	return TitleStyle.Render("PARENT MODE") + "\n\n" +
		SuccessStyle.Render("Access granted!") + "\n\n" +
		"Q - Quit game\n" +
		"ESC - Return to game"
}
