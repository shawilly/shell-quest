package shell_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

func TestExecutor_UnknownCommand(t *testing.T) {
	fs := shell.NewFS()
	ex := shell.NewExecutor(fs)
	result := ex.Execute("notacommand", "/")
	if result.Error == "" {
		t.Error("expected error for unknown command")
	}
}

func TestExecutor_EmptyInput(t *testing.T) {
	fs := shell.NewFS()
	ex := shell.NewExecutor(fs)
	result := ex.Execute("", "/")
	if result.Error != "" || result.Output != "" {
		t.Errorf("empty input should produce empty result: %+v", result)
	}
}
