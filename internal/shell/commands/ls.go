package commands

import (
	"strings"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

type Ls struct{}

func NewLs() *Ls { return &Ls{} }

func (l *Ls) Name() string      { return "ls" }
func (l *Ls) Aliases() []string { return nil }
func (l *Ls) Tier() string      { return "beginner" }

func (l *Ls) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	showAll := false
	target := cwd
	for _, arg := range args {
		if arg == "-a" {
			showAll = true
		} else if !strings.HasPrefix(arg, "-") {
			target = shell.ResolvePath(cwd, arg)
		}
	}

	var entries []*shell.Node
	var err error
	if showAll {
		entries, err = fs.ListDirAll(target)
	} else {
		entries, err = fs.ListDir(target)
	}
	if err != nil {
		return shell.Result{Error: "ls: " + err.Error()}
	}

	if len(entries) == 0 {
		return shell.Result{Output: ""}
	}

	var parts []string
	for _, e := range entries {
		name := e.Name
		if e.Type == shell.NodeDir {
			name += "/"
		}
		parts = append(parts, name)
	}
	return shell.Result{Output: strings.Join(parts, "  "), Event: &shell.Event{Type: "ls", Path: target}}
}
