package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestMv_MovesFile(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.WriteFile("/src.txt", "hi", false)
	cmd := commands.NewMv()
	result := cmd.Run([]string{"src.txt", "dst.txt"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if _, err := fs.Stat("/src.txt"); err == nil {
		t.Error("source should be gone")
	}
	if _, err := fs.Stat("/dst.txt"); err != nil {
		t.Error("destination should exist")
	}
}
