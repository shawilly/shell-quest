package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestMkdir_CreatesDir(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewMkdir()
	result := cmd.Run([]string{"newdir"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	node, err := fs.Stat("/newdir")
	if err != nil || node.Type != shell.NodeDir {
		t.Errorf("expected dir at /newdir: %v", err)
	}
}

func TestMkdir_NoArgs_Errors(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewMkdir()
	result := cmd.Run(nil, "/", fs)
	if result.Error == "" {
		t.Error("expected error with no args")
	}
}
