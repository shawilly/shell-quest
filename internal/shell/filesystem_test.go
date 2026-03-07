package shell_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

func TestNewFS_HasRoot(t *testing.T) {
	fs := shell.NewFS()
	node, err := fs.Stat("/")
	if err != nil {
		t.Fatal(err)
	}
	if node.Name != "/" || node.Type != shell.NodeDir {
		t.Errorf("unexpected root: %+v", node)
	}
}

func TestMkdir_AndStat(t *testing.T) {
	fs := shell.NewFS()
	if err := fs.Mkdir("/island", false); err != nil {
		t.Fatal(err)
	}
	node, err := fs.Stat("/island")
	if err != nil {
		t.Fatal(err)
	}
	if node.Type != shell.NodeDir {
		t.Errorf("expected dir, got %v", node.Type)
	}
}

func TestStat_MissingPath_Errors(t *testing.T) {
	fs := shell.NewFS()
	_, err := fs.Stat("/nonexistent")
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestListDir(t *testing.T) {
	fs := shell.NewFS()
	fs.Mkdir("/island", false)
	fs.Mkdir("/island/cave", false)
	fs.WriteFile("/island/note.txt", "hello", false)

	entries, err := fs.ListDir("/island")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2, got %d: %v", len(entries), entries)
	}
}

func TestListDir_HiddenFilesNotShownByDefault(t *testing.T) {
	fs := shell.NewFS()
	fs.Mkdir("/island", false)
	fs.WriteFile("/island/.secret", "shh", true)
	fs.WriteFile("/island/visible.txt", "hi", false)

	entries, err := fs.ListDir("/island")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name != "visible.txt" {
		t.Errorf("expected only visible.txt, got %v", entries)
	}
}

func TestListDirAll_ShowsHidden(t *testing.T) {
	fs := shell.NewFS()
	fs.Mkdir("/island", false)
	fs.WriteFile("/island/.secret", "shh", true)
	fs.WriteFile("/island/visible.txt", "hi", false)

	entries, err := fs.ListDirAll("/island")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries including hidden, got %d", len(entries))
	}
}
