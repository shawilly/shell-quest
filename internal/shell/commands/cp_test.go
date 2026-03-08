package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestCp_CopiesFile(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.WriteFile("/src.txt", "hello", false)
	cmd := commands.NewCp()
	result := cmd.Run([]string{"src.txt", "dst.txt"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	node, err := fs.Stat("/dst.txt")
	if err != nil || node.Content != "hello" {
		t.Errorf("expected dst.txt with 'hello': %v", err)
	}
}

func TestCp_MissingArgs_Errors(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewCp()
	result := cmd.Run([]string{"only-one"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error with only one arg")
	}
}
