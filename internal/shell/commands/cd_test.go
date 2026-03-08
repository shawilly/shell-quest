package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestCd_ValidDir(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.Mkdir("/island/cave", false)

	cmd := commands.NewCd()
	result := cmd.Run([]string{"cave"}, "/island", fs)
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if result.NewCWD != "/island/cave" {
		t.Errorf("expected /island/cave, got %q", result.NewCWD)
	}
}

func TestCd_DotDot(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	cmd := commands.NewCd()
	result := cmd.Run([]string{".."}, "/island", fs)
	if result.NewCWD != "/" {
		t.Errorf("expected /, got %q", result.NewCWD)
	}
}

func TestCd_NonexistentDir_Errors(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewCd()
	result := cmd.Run([]string{"nope"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error for nonexistent dir")
	}
}

func TestCd_FileNotDir_Errors(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.WriteFile("/note.txt", "hi", false)
	cmd := commands.NewCd()
	result := cmd.Run([]string{"note.txt"}, "/", fs)
	if result.Error == "" {
		t.Error("expected error when cd-ing into a file")
	}
}
