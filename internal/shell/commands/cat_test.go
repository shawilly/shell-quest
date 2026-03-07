package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestCat_ReadsFile(t *testing.T) {
	fs := shell.NewFS()
	fs.WriteFile("/note.txt", "treasure here!", false)
	cmd := commands.NewCat()
	result := cmd.Run([]string{"note.txt"}, "/", fs)
	if result.Output != "treasure here!" {
		t.Errorf("got %q", result.Output)
	}
}

func TestCat_EmitsEvent(t *testing.T) {
	fs := shell.NewFS()
	fs.WriteFile("/clue.txt", "go north", false)
	cmd := commands.NewCat()
	result := cmd.Run([]string{"clue.txt"}, "/", fs)
	if result.Event == nil || result.Event.Type != "cat" {
		t.Error("expected cat event")
	}
}

func TestCat_Directory_Errors(t *testing.T) {
	fs := shell.NewFS()
	fs.Mkdir("/island", false)
	cmd := commands.NewCat()
	result := cmd.Run([]string{"island"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error for directory")
	}
}
