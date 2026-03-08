package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

type Pwd struct{}

func NewPwd() *Pwd { return &Pwd{} }

func (p *Pwd) Name() string      { return "pwd" }
func (p *Pwd) Aliases() []string { return nil }
func (p *Pwd) Tier() string      { return "beginner" }

func (p *Pwd) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	return shell.Result{Output: cwd, Event: &shell.Event{Type: "pwd", Path: cwd}}
}
