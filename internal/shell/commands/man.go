package commands

import (
	"fmt"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

var manPages = map[string]string{
	"ls":      "ls - list directory contents\n\nUsage: ls [-a] [directory]\n  -a   show hidden files too (files starting with .)",
	"cd":      "cd - change directory\n\nUsage: cd [directory]\n  Without arguments, goes to /",
	"pwd":     "pwd - print working directory\n\nUsage: pwd\n  Shows your current location.",
	"cat":     "cat - read and print a file\n\nUsage: cat <file>",
	"echo":    "echo - print a message\n\nUsage: echo <message>",
	"mkdir":   "mkdir - make a directory\n\nUsage: mkdir <dirname>",
	"touch":   "touch - create an empty file\n\nUsage: touch <filename>",
	"rm":      "rm - remove a file\n\nUsage: rm <file>\n  Cannot remove directories.",
	"cp":      "cp - copy a file\n\nUsage: cp <source> <destination>",
	"mv":      "mv - move or rename a file\n\nUsage: mv <source> <destination>",
	"find":    "find - search for files\n\nUsage: find [dir] -name <pattern>",
	"grep":    "grep - search inside files\n\nUsage: grep <pattern> <file>\nPrints lines that contain the pattern.",
	"chmod":   "chmod - change file permissions\n\nUsage: chmod <mode> <file>\nExample: chmod 755 myfile.txt",
	"history": "history - show past commands\n\nUsage: history\n  Lists all commands you have typed.",
	"pipe":    "pipe ( | ) - connect commands\n\nExample: cat file.txt | grep treasure\nSends output of one command to the next.",
}

type Man struct{}

func NewMan() *Man { return &Man{} }

func (m *Man) Name() string      { return "man" }
func (m *Man) Aliases() []string { return nil }
func (m *Man) Tier() string      { return "master" }

func (m *Man) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	if len(args) == 0 {
		return shell.Result{Error: "man: what manual page do you want?"}
	}
	page, ok := manPages[args[0]]
	if !ok {
		return shell.Result{Error: fmt.Sprintf("man: no manual entry for %s", args[0])}
	}
	return shell.Result{Output: page}
}
