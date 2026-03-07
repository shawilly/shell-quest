package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestEcho(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewEcho()
	result := cmd.Run([]string{"hello", "pirate"}, "/", fs)
	if result.Output != "hello pirate" {
		t.Errorf("got %q", result.Output)
	}
}

func TestEcho_NoArgs(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewEcho()
	result := cmd.Run(nil, "/", fs)
	if result.Output != "" {
		t.Errorf("expected empty output, got %q", result.Output)
	}
}
