package commands_test

import (
	"strings"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestFind_FindsFileByName(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.Mkdir("/island/cave", false)
	_ = fs.WriteFile("/island/cave/treasure.txt", "gold!", false)

	cmd := commands.NewFind()
	result := cmd.Run([]string{"-name", "treasure.txt"}, "/island", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if !strings.Contains(result.Output, "treasure.txt") {
		t.Errorf("expected treasure.txt in output: %q", result.Output)
	}
}

func TestFind_NoMatch_ReturnsNothingFound(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)

	cmd := commands.NewFind()
	result := cmd.Run([]string{"-name", "missing.txt"}, "/island", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if !strings.Contains(result.Output, "nothing found") {
		t.Errorf("expected 'nothing found', got: %q", result.Output)
	}
}
