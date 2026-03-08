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
