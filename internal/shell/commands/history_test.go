package commands_test

import (
	"strings"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestHistory_Empty(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewHistory(func() []string { return nil })
	result := cmd.Run(nil, "/", fs)
	if !strings.Contains(result.Output, "no history") {
		t.Errorf("expected 'no history': %q", result.Output)
	}
}

func TestHistory_ShowsPastCommands(t *testing.T) {
	fs := shell.NewFS()
	hist := []string{"ls", "cd island", "cat note.txt"}
	cmd := commands.NewHistory(func() []string { return hist })
	result := cmd.Run(nil, "/", fs)
	if !strings.Contains(result.Output, "ls") || !strings.Contains(result.Output, "cat note.txt") {
		t.Errorf("expected history entries in output: %q", result.Output)
	}
}
