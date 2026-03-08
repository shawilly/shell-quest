package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Exit struct{}

func NewExit() *Exit { return &Exit{} }

func (e *Exit) Name() string      { return "exit" }
func (e *Exit) Aliases() []string { return []string{"quit"} }
func (e *Exit) Tier() string      { return "beginner" }

func (e *Exit) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	return shell.Result{ExitLevel: true}
}
