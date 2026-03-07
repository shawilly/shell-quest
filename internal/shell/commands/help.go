package commands

import "github.com/shanewilliams/shell-quest/internal/shell"

var helpText = map[string]string{
	"beginner": `Available Commands:
  ls        - list what's around you
  cd <dir>  - move into a directory
  pwd       - show where you are
  cat <file>- read a file
  echo <msg>- print a message
  clear     - clear the screen
  help      - show this message`,

	"explorer": `Available Commands (Explorer):
  [all beginner commands, plus...]
  mkdir <dir>    - create a new directory
  touch <file>   - create an empty file
  cp <src> <dst> - copy a file
  mv <src> <dst> - move/rename a file
  rm <file>      - delete a file
  find <name>    - search for a file`,

	"master": `Available Commands (Master):
  [all explorer commands, plus...]
  grep <pattern> <file> - search inside files
  chmod <mode> <file>   - change permissions
  man <command>         - read the manual
  history               - show past commands
  cmd1 | cmd2           - pipe commands
  *.txt                 - glob patterns
  cmd > file            - redirect output`,
}

type Help struct {
	tier string
}

func NewHelp(tier string) *Help { return &Help{tier: tier} }

func (h *Help) Name() string      { return "help" }
func (h *Help) Aliases() []string { return []string{"?"} }
func (h *Help) Tier() string      { return "beginner" }

func (h *Help) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	text, ok := helpText[h.tier]
	if !ok {
		text = helpText["beginner"]
	}
	return shell.Result{Output: text}
}
