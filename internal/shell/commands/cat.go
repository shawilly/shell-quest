package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Cat struct{}

func NewCat() *Cat { return &Cat{} }

func (c *Cat) Name() string      { return "cat" }
func (c *Cat) Aliases() []string { return nil }
func (c *Cat) Tier() string      { return "beginner" }

func (c *Cat) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) == 0 {
		return shell.Result{Error: "cat: missing file operand"}
	}
	p := shell.ResolvePath(cwd, args[0])
	node, err := fs.Stat(p)
	if err != nil {
		return shell.Result{Error: "cat: " + err.Error()}
	}
	if node.Type == shell.NodeDir {
		return shell.Result{Error: "cat: " + args[0] + ": Is a directory"}
	}
	return shell.Result{Output: node.Content, Event: &shell.Event{Type: "cat", Path: p}}
}
