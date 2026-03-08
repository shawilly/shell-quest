package commands_test

import (
	"strings"
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestLs_ListsCurrentDir(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.Mkdir("/island/cave", false)
	_ = fs.WriteFile("/island/note.txt", "hello", false)

	cmd := commands.NewLs()
	result := cmd.Run(nil, "/island", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if !strings.Contains(result.Output, "cave") || !strings.Contains(result.Output, "note.txt") {
		t.Errorf("expected cave and note.txt in output: %q", result.Output)
	}
}

func TestLs_HidesHiddenFiles(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.WriteFile("/island/.secret", "shh", true)
	_ = fs.WriteFile("/island/visible.txt", "hi", false)

	cmd := commands.NewLs()
	result := cmd.Run(nil, "/island", fs)
	if strings.Contains(result.Output, ".secret") {
		t.Errorf("should not show hidden file without -a: %q", result.Output)
	}
}

func TestLs_WithDashA_ShowsHidden(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.WriteFile("/island/.secret", "shh", true)

	cmd := commands.NewLs()
	result := cmd.Run([]string{"-a"}, "/island", fs)
	if !strings.Contains(result.Output, ".secret") {
		t.Errorf("expected .secret with -a: %q", result.Output)
	}
}

func TestLs_NonexistentDir_ReturnsError(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewLs()
	result := cmd.Run([]string{"/nope"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error for nonexistent dir")
	}
}
