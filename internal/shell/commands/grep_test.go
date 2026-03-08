package commands_test

import (
	"strings"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestGrep_FindsMatch(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.WriteFile("/clue.txt", "line one\ntreasure is here\nline three", false)
	cmd := commands.NewGrep()
	result := cmd.Run([]string{"treasure", "clue.txt"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if !strings.Contains(result.Output, "treasure is here") {
		t.Errorf("expected match in output: %q", result.Output)
	}
}

func TestGrep_NoMatch(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.WriteFile("/file.txt", "nothing here", false)
	cmd := commands.NewGrep()
	result := cmd.Run([]string{"missing", "file.txt"}, "/", fs)
	if !strings.Contains(result.Output, "no matches") {
		t.Errorf("expected 'no matches': %q", result.Output)
	}
}

func TestGrep_MissingArgs_Errors(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewGrep()
	result := cmd.Run([]string{"pattern"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error with only one arg")
	}
}
