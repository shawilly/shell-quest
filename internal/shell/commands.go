package shell

// Result is the output of a command execution.
type Result struct {
	Output    string
	NewCWD    string
	Error     string
	Event     *Event
	Clear     bool
	ExitLevel bool
}

// Event is emitted when a command triggers a game narrative event.
type Event struct {
	Type string
	Path string
}

// Commander is the interface all shell commands implement.
type Commander interface {
	Name() string
	Aliases() []string
	Tier() string
	Run(args []string, cwd string, fs *FS) Result
}
