package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Rm struct{}

func NewRm() *Rm { return &Rm{} }

func (r *Rm) Name() string      { return "rm" }
func (r *Rm) Aliases() []string { return nil }
func (r *Rm) Tier() string      { return "explorer" }

func (r *Rm) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) == 0 {
		return shell.Result{Error: "rm: missing operand"}
	}
	p := shell.ResolvePath(cwd, args[0])
	node, err := fs.Stat(p)
	if err != nil {
		return shell.Result{Error: "rm: " + err.Error()}
	}
	if node.Type == shell.NodeDir {
		return shell.Result{Error: "rm: cannot remove directory (use rmdir)"}
	}
	if err := fs.Remove(p); err != nil {
		return shell.Result{Error: "rm: " + err.Error()}
	}
	return shell.Result{}
}
