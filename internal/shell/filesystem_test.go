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
	if err := fs.Mkdir("/island", false); err != nil {
		t.Fatal(err)
	}
	if err := fs.Mkdir("/island/cave", false); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile("/island/note.txt", "hello", false); err != nil {
		t.Fatal(err)
	}

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
	if err := fs.Mkdir("/island", false); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile("/island/.secret", "shh", true); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile("/island/visible.txt", "hi", false); err != nil {
		t.Fatal(err)
	}

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
	_ = fs.Mkdir("/island", false)
	_ = fs.WriteFile("/island/.secret", "shh", true)
	_ = fs.WriteFile("/island/visible.txt", "hi", false)

	entries, err := fs.ListDirAll("/island")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries including hidden, got %d", len(entries))
	}
}

func TestRemove_File(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.WriteFile("/island/note.txt", "hi", false)
	if err := fs.Remove("/island/note.txt"); err != nil {
		t.Fatal(err)
	}
	_, err := fs.Stat("/island/note.txt")
	if err == nil {
		t.Error("expected error after removal")
	}
}

func TestCopy_File(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/a", false)
	_ = fs.Mkdir("/b", false)
	_ = fs.WriteFile("/a/file.txt", "content", false)
	if err := fs.Copy("/a/file.txt", "/b/file.txt"); err != nil {
		t.Fatal(err)
	}
	node, err := fs.Stat("/b/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	if node.Content != "content" {
		t.Errorf("expected 'content', got %q", node.Content)
	}
}

func TestMove_File(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/a", false)
	_ = fs.Mkdir("/b", false)
	_ = fs.WriteFile("/a/file.txt", "content", false)
	if err := fs.Move("/a/file.txt", "/b/file.txt"); err != nil {
		t.Fatal(err)
	}
	if _, err := fs.Stat("/a/file.txt"); err == nil {
		t.Error("expected source to be gone")
	}
	if _, err := fs.Stat("/b/file.txt"); err != nil {
		t.Error("expected destination to exist")
	}
}

func TestFS_Serialize_Deserialize(t *testing.T) {
	fs := shell.NewFS()
	_ = fs.Mkdir("/island", false)
	_ = fs.Mkdir("/island/cave", false)
	_ = fs.WriteFile("/island/cave/note.txt", "hello pirate", false)
	_ = fs.WriteFile("/island/cave/.secret", "hidden treasure", true)

	json, err := fs.Serialize()
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	if json == "" {
		t.Fatal("Serialize returned empty string")
	}

	// Deserialize into a new FS
	fs2, err := shell.DeserializeFS(json)
	if err != nil {
		t.Fatalf("DeserializeFS: %v", err)
	}

	// Check the note file
	node, err := fs2.Stat("/island/cave/note.txt")
	if err != nil {
		t.Fatalf("Stat after deserialize: %v", err)
	}
	if node.Content != "hello pirate" {
		t.Errorf("expected 'hello pirate', got %q", node.Content)
	}

	// Check hidden file is preserved
	node2, err := fs2.Stat("/island/cave/.secret")
	if err != nil {
		t.Fatalf("Stat hidden: %v", err)
	}
	if !node2.Hidden {
		t.Error("expected .secret to be hidden")
	}
}
