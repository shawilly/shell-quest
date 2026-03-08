package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func newNameInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "e.g. Blackbeard"
	ti.CharLimit = 20
	ti.Width = 24
	ti.Prompt = "> "
	ti.PromptStyle = PromptStyle
	ti.TextStyle = OutputStyle
	return ti
}

func (m Model) nameInputView() string {
	return TitleStyle.Render("SHELL QUEST — Enter Your Pirate Name") + "\n\n" +
		"What is your name, young pirate?\n\n" +
		m.nameInput.View() + "\n\n" +
		"Press Enter when done."
}
