package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Mkdir struct{}

func NewMkdir() *Mkdir { return &Mkdir{} }

func (m *Mkdir) Name() string      { return "mkdir" }
func (m *Mkdir) Aliases() []string { return nil }
func (m *Mkdir) Tier() string      { return "explorer" }

func (m *Mkdir) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) == 0 {
		return shell.Result{Error: "mkdir: missing operand"}
	}
	p := shell.ResolvePath(cwd, args[0])
	if err := fs.Mkdir(p, false); err != nil {
		return shell.Result{Error: "mkdir: " + err.Error()}
	}
	return shell.Result{Event: &shell.Event{Type: "mkdir", Path: p}}
}
