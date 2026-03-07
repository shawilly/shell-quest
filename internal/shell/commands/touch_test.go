package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestTouch_CreatesFile(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewTouch()
	result := cmd.Run([]string{"newfile.txt"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	node, err := fs.Stat("/newfile.txt")
	if err != nil || node.Type != shell.NodeFile {
		t.Errorf("expected file at /newfile.txt: %v", err)
	}
}

func TestTouch_ExistingFile_NoOp(t *testing.T) {
	fs := shell.NewFS()
	fs.WriteFile("/existing.txt", "content", false)
	cmd := commands.NewTouch()
	result := cmd.Run([]string{"existing.txt"}, "/", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	// Content should be unchanged
	node, _ := fs.Stat("/existing.txt")
	if node.Content != "content" {
		t.Errorf("touch modified existing file content: %q", node.Content)
	}
}
