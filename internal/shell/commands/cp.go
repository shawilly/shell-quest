package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Cp struct{}

func NewCp() *Cp { return &Cp{} }

func (c *Cp) Name() string      { return "cp" }
func (c *Cp) Aliases() []string { return nil }
func (c *Cp) Tier() string      { return "explorer" }

func (c *Cp) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) < 2 {
		return shell.Result{Error: "cp: missing operands"}
	}
	src := shell.ResolvePath(cwd, args[0])
	dst := shell.ResolvePath(cwd, args[1])
	if err := fs.Copy(src, dst); err != nil {
		return shell.Result{Error: "cp: " + err.Error()}
	}
	return shell.Result{}
}
