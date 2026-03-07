package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Clear struct{}

func NewClear() *Clear { return &Clear{} }

func (c *Clear) Name() string      { return "clear" }
func (c *Clear) Aliases() []string { return []string{"cls"} }
func (c *Clear) Tier() string      { return "beginner" }

func (c *Clear) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	return shell.Result{Clear: true}
}
