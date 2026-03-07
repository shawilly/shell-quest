package shell

import (
	"path"
	"strings"
)

// ParsedCommand represents one command in a pipeline.
type ParsedCommand struct {
	Name       string
	Args       []string
	RedirectTo string
}

// Parse tokenizes input into a pipeline of commands.
// Handles: pipes (|), redirect (>), glob expansion (*).
func Parse(input, cwd string, fs *FS) []ParsedCommand {
	segments := strings.Split(input, "|")
	var pipeline []ParsedCommand
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		// Check for redirect
		redirect := ""
		if idx := strings.Index(seg, ">"); idx != -1 {
			redirect = strings.TrimSpace(seg[idx+1:])
			seg = strings.TrimSpace(seg[:idx])
		}
		parts := strings.Fields(seg)
		if len(parts) == 0 {
			continue
		}
		// Expand globs in args
		var expanded []string
		for _, part := range parts[1:] {
			if strings.Contains(part, "*") {
				matches := expandGlob(cwd, part, fs)
				expanded = append(expanded, matches...)
			} else {
				expanded = append(expanded, part)
			}
		}
		pipeline = append(pipeline, ParsedCommand{
			Name:       parts[0],
			Args:       expanded,
			RedirectTo: strings.TrimSpace(redirect),
		})
	}
	return pipeline
}

func expandGlob(cwd, pattern string, fs *FS) []string {
	dir := cwd
	base := pattern
	if strings.Contains(pattern, "/") {
		dir = ResolvePath(cwd, path.Dir(pattern))
		base = path.Base(pattern)
	}
	entries, err := fs.ListDir(dir)
	if err != nil {
		return []string{pattern}
	}
	starIdx := strings.Index(base, "*")
	prefix := base[:starIdx]
	suffix := base[starIdx+1:]
	var matches []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name, prefix) && strings.HasSuffix(e.Name, suffix) {
			if dir == cwd {
				matches = append(matches, e.Name)
			} else {
				matches = append(matches, path.Join(dir, e.Name))
			}
		}
	}
	if len(matches) == 0 {
		return []string{pattern}
	}
	return matches
}
