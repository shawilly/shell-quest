package commands

import (
	"fmt"
	"strings"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

type Grep struct{}

func NewGrep() *Grep { return &Grep{} }

func (g *Grep) Name() string      { return "grep" }
func (g *Grep) Aliases() []string { return nil }
func (g *Grep) Tier() string      { return "master" }

func (g *Grep) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) < 2 {
		return shell.Result{Error: "grep: usage: grep <pattern> <file>"}
	}
	pattern := args[0]
	p := shell.ResolvePath(cwd, args[1])
	node, err := fs.Stat(p)
	if err != nil {
		return shell.Result{Error: "grep: " + err.Error()}
	}
	if node.Type == shell.NodeDir {
		return shell.Result{Error: "grep: " + args[1] + ": Is a directory"}
	}
	var matches []string
	for i, line := range strings.Split(node.Content, "\n") {
		if strings.Contains(line, pattern) {
			matches = append(matches, fmt.Sprintf("%d: %s", i+1, line))
		}
	}
	if len(matches) == 0 {
		return shell.Result{Output: "(no matches)"}
	}
	return shell.Result{Output: strings.Join(matches, "\n")}
}
