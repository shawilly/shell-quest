package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Chmod struct{}

func NewChmod() *Chmod { return &Chmod{} }

func (c *Chmod) Name() string      { return "chmod" }
func (c *Chmod) Aliases() []string { return nil }
func (c *Chmod) Tier() string      { return "master" }

func (c *Chmod) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) < 2 {
		return shell.Result{Error: "chmod: usage: chmod <mode> <file>"}
	}
	mode := args[0]
	p := shell.ResolvePath(cwd, args[1])
	node, err := fs.Stat(p)
	if err != nil {
		return shell.Result{Error: "chmod: " + err.Error()}
	}
	node.Permissions = mode
	return shell.Result{Output: "Changed permissions of " + args[1] + " to " + mode}
}
