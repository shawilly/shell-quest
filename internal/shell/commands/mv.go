package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Mv struct{}

func NewMv() *Mv { return &Mv{} }

func (m *Mv) Name() string      { return "mv" }
func (m *Mv) Aliases() []string { return nil }
func (m *Mv) Tier() string      { return "explorer" }

func (m *Mv) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) < 2 {
		return shell.Result{Error: "mv: missing operands"}
	}
	src := shell.ResolvePath(cwd, args[0])
	dst := shell.ResolvePath(cwd, args[1])
	if err := fs.Move(src, dst); err != nil {
		return shell.Result{Error: "mv: " + err.Error()}
	}
	return shell.Result{Event: &shell.Event{Type: "mv", Path: dst}}
}
