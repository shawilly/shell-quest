package commands_test

import (
	"testing"

	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
)

func TestPwd(t *testing.T) {
	fs := shell.NewFS()
	cmd := commands.NewPwd()
	result := cmd.Run(nil, "/island/cave", fs)
	if result.Output != "/island/cave" {
		t.Errorf("expected '/island/cave', got %q", result.Output)
	}
	if result.Error != "" {
		t.Errorf("unexpected error: %s", result.Error)
	}
}
