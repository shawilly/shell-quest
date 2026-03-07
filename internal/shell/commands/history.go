package commands

import (
	"fmt"
	"strings"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

type History struct {
	getHistory func() []string
}

func NewHistory(getHistory func() []string) *History {
	return &History{getHistory: getHistory}
}

func (h *History) Name() string      { return "history" }
func (h *History) Aliases() []string { return nil }
func (h *History) Tier() string      { return "master" }

func (h *History) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	hist := h.getHistory()
	if len(hist) == 0 {
		return shell.Result{Output: "(no history yet)"}
	}
	var lines []string
	for i, cmd := range hist {
		lines = append(lines, fmt.Sprintf("%3d  %s", i+1, cmd))
	}
	return shell.Result{Output: strings.Join(lines, "\n")}
}
