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
	parts := strings.Fields(input)
	name, args := parts[0], parts[1:]
	cmd, ok := e.registry[name]
	if !ok {
		return Result{Error: name + ": command not found"}
	}
	return cmd.Run(args, cwd, e.fs)
}
