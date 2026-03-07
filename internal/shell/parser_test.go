package shell_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

func TestParse_SimpleCommand(t *testing.T) {
	fs := shell.NewFS()
	pipeline := shell.Parse("ls -a", "/", fs)
	if len(pipeline) != 1 {
		t.Fatalf("expected 1 command, got %d", len(pipeline))
	}
	if pipeline[0].Name != "ls" || len(pipeline[0].Args) != 1 || pipeline[0].Args[0] != "-a" {
		t.Errorf("unexpected parse: %+v", pipeline[0])
	}
}

func TestParse_Pipe(t *testing.T) {
	fs := shell.NewFS()
	pipeline := shell.Parse("cat file.txt | grep treasure", "/", fs)
	if len(pipeline) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(pipeline))
	}
	if pipeline[0].Name != "cat" || pipeline[1].Name != "grep" {
		t.Errorf("unexpected pipeline: %+v", pipeline)
	}
}

func TestParse_Redirect(t *testing.T) {
	fs := shell.NewFS()
	pipeline := shell.Parse("echo hello > out.txt", "/", fs)
	if len(pipeline) != 1 {
		t.Fatalf("expected 1 command, got %d", len(pipeline))
	}
	if pipeline[0].RedirectTo != "out.txt" {
		t.Errorf("expected redirect to out.txt: %+v", pipeline[0])
	}
}

func TestParse_Glob_Expands(t *testing.T) {
	fs := shell.NewFS()
	fs.WriteFile("/island/cave.txt", "", false)
	fs.WriteFile("/island/note.txt", "", false)
	fs.WriteFile("/island/readme.md", "", false)
	pipeline := shell.Parse("cat *.txt", "/island", fs)
	if len(pipeline) != 1 {
		t.Fatalf("expected 1 command, got %d", len(pipeline))
	}
	// *.txt should expand to cave.txt and note.txt (sorted)
	if len(pipeline[0].Args) != 2 {
		t.Errorf("expected 2 args from glob, got %d: %v", len(pipeline[0].Args), pipeline[0].Args)
	}
}
