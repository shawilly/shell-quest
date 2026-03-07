package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestRm_DeletesFile(t *testing.T) {
	fs := shell.NewFS()
	fs.WriteFile("/gone.txt", "bye", false)
	cmd := commands.NewRm()
	result := cmd.Run([]string{"gone.txt"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if _, err := fs.Stat("/gone.txt"); err == nil {
		t.Error("file should be gone")
	}
}

func TestRm_Directory_Errors(t *testing.T) {
	fs := shell.NewFS()
	fs.Mkdir("/dir", false)
	cmd := commands.NewRm()
	result := cmd.Run([]string{"dir"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error when removing directory")
	}
}
