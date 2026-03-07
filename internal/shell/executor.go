package shell

import "strings"

type Executor struct {
	fs       *FS
	registry map[string]Commander
	History  []string
}

func NewExecutor(fs *FS) *Executor {
	return &Executor{
		fs:       fs,
		registry: make(map[string]Commander),
	}
}

func (e *Executor) Register(cmd Commander) {
	e.registry[cmd.Name()] = cmd
	for _, alias := range cmd.Aliases() {
		e.registry[alias] = cmd
	}
}

func (e *Executor) Execute(input, cwd string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{}
	}
	e.History = append(e.History, input)

	pipeline := Parse(input, cwd, e.fs)
	if len(pipeline) == 0 {
		return Result{}
	}

	var lastOutput string
	var lastResult Result

	for i, pc := range pipeline {
		cmd, ok := e.registry[pc.Name]
		if !ok {
			return Result{Error: pc.Name + ": command not found"}
		}
		args := pc.Args
		// For pipes: append previous output as extra arg (simplified pipe)
		if i > 0 && lastOutput != "" {
			args = append(args, lastOutput)
		}
		lastResult = cmd.Run(args, cwd, e.fs)
		if lastResult.Error != "" {
			return lastResult
		}
		// Handle redirect
		if pc.RedirectTo != "" {
			p := ResolvePath(cwd, pc.RedirectTo)
			e.fs.WriteFile(p, lastResult.Output, false)
			lastResult.Output = ""
		}
		lastOutput = lastResult.Output
	}
	return lastResult
}
