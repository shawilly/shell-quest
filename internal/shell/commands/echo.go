package commands

import (
	"strings"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

type Echo struct{}

func NewEcho() *Echo { return &Echo{} }

func (e *Echo) Name() string      { return "echo" }
func (e *Echo) Aliases() []string { return nil }
func (e *Echo) Tier() string      { return "beginner" }

func (e *Echo) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	return shell.Result{Output: strings.Join(args, " ")}
}
