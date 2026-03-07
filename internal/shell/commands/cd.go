package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Cd struct{}

func NewCd() *Cd { return &Cd{} }

func (c *Cd) Name() string      { return "cd" }
func (c *Cd) Aliases() []string { return nil }
func (c *Cd) Tier() string      { return "beginner" }

func (c *Cd) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	target := "/"
	if len(args) > 0 {
		target = shell.ResolvePath(cwd, args[0])
	}
	node, err := fs.Stat(target)
	if err != nil {
		return shell.Result{Error: "cd: " + err.Error()}
	}
	if node.Type != shell.NodeDir {
		return shell.Result{Error: "cd: not a directory: " + args[0]}
	}
	return shell.Result{NewCWD: target, Event: &shell.Event{Type: "cd", Path: target}}
}
