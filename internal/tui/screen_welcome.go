package tui

import (
	"strings"

	"github.com/shanewilliams/shell-quest/content"
)

func (m Model) welcomeView() string {
	var b strings.Builder
	b.WriteString(content.PirateArt + "\n")
	b.WriteString(TitleStyle.Render("Welcome to Shell Quest!") + "\n\n")
	b.WriteString("Learn the secrets of the command line\n")
	b.WriteString("through pirate treasure hunting!\n\n")
	b.WriteString("Press ENTER to begin your adventure!")
	return b.String()
}
