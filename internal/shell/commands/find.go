package commands

import (
	"path"
	"strings"

	"github.com/shanewilliams/shell-quest/internal/shell"
)

type Find struct{}

func NewFind() *Find { return &Find{} }

func (f *Find) Name() string      { return "find" }
func (f *Find) Aliases() []string { return nil }
func (f *Find) Tier() string      { return "explorer" }

func (f *Find) Run(args []string, cwd string, fs *shell.FS) shell.Result {
	searchDir := cwd
	pattern := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "-name" && i+1 < len(args) {
			pattern = args[i+1]
			i++
		} else {
			searchDir = shell.ResolvePath(cwd, args[i])
		}
	}
	var results []string
	var walk func(dir string)
	walk = func(dir string) {
		entries, err := fs.ListDirAll(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			p := path.Join(dir, e.Name)
			if pattern == "" || strings.Contains(e.Name, pattern) {
				results = append(results, p)
			}
			if e.Type == shell.NodeDir {
				walk(p)
			}
		}
	}
	walk(searchDir)
	if len(results) == 0 {
		return shell.Result{Output: "(nothing found)"}
	}
	return shell.Result{Output: strings.Join(results, "\n"), Event: &shell.Event{Type: "find", Path: searchDir}}
}
