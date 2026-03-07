package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Touch struct{}

func NewTouch() *Touch { return &Touch{} }

func (t *Touch) Name() string      { return "touch" }
func (t *Touch) Aliases() []string { return nil }
func (t *Touch) Tier() string      { return "explorer" }

func (t *Touch) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) == 0 {
		return shell.Result{Error: "touch: missing operand"}
	}
	p := shell.ResolvePath(cwd, args[0])
	if _, err := fs.Stat(p); err != nil {
		if err := fs.WriteFile(p, "", false); err != nil {
			return shell.Result{Error: "touch: " + err.Error()}
		}
	}
	return shell.Result{Event: &shell.Event{Type: "touch", Path: p}}
}
