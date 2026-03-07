package commands_test

import (
	"strings"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestMan_KnownCommand(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewMan()
	result := cmd.Run([]string{"ls"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if !strings.Contains(result.Output, "ls") {
		t.Errorf("expected ls info in output: %q", result.Output)
	}
}

func TestMan_UnknownCommand_Errors(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewMan()
	result := cmd.Run([]string{"notacommand"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error for unknown man page")
	}
}

func TestMan_NoArgs_Errors(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewMan()
	result := cmd.Run(nil, "/", fs)
	if result.Error == "" {
		t.Error("expected error with no args")
	}
}
