package tui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shanewilliams/shell-quest/internal/db"
	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
	"github.com/shanewilliams/shell-quest/internal/world"
)

type AppState int

const (
	StateGame AppState = iota
	StateProfileSelect
	StateTierSelect
	StateNameInput
)

// Model is the root Bubble Tea model.
type Model struct {
	state    AppState
	db       *db.DB
	player   *db.Player
	fs       *shell.FS
	executor *shell.Executor
	runner   *world.MissionRunner
	cwd      string
	width    int
	height   int

	// shell pane
	outputLines []string
	inputBuf    string
	maxLines    int

	// story pane
	clueText  string
	storyText string
	failCount int

	// profile selection
	profiles    []*db.Player
	selectedIdx int
	nameInput   string
}

// NewGameModel creates a model ready to play a mission.
func NewGameModel(d *db.DB, player *db.Player, fs *shell.FS, ex *shell.Executor, runner *world.MissionRunner, startCWD, clue string) Model {
	return Model{
		state:     StateGame,
		db:        d,
		player:    player,
		fs:        fs,
		executor:  ex,
		runner:    runner,
		cwd:       startCWD,
		clueText:  clue,
		storyText: "Welcome to Skull Island, young pirate! Arr!",
		maxLines:  20,
	}
}

// NewStartupModel creates a model starting at the profile select screen.
func NewStartupModel(d *db.DB) Model {
	players, _ := d.ListPlayers()
	return Model{
		state:    StateProfileSelect,
		db:       d,
		profiles: players,
		maxLines: 20,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch m.state {
		case StateGame:
			return m.handleKey(msg)
		case StateProfileSelect:
			return m.handleProfileKey(msg)
		case StateTierSelect:
			return m.handleTierKey(msg)
		case StateNameInput:
			return m.handleNameKey(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case StateGame:
		return m.gameView()
	case StateProfileSelect:
		return m.profileSelectView()
	case StateTierSelect:
		return m.tierSelectView()
	case StateNameInput:
		return m.nameInputView()
	default:
		return "Loading..."
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		return m.submitCommand()
	case tea.KeyBackspace:
		if len(m.inputBuf) > 0 {
			m.inputBuf = m.inputBuf[:len(m.inputBuf)-1]
		}
	default:
		if msg.Type == tea.KeyRunes {
			m.inputBuf += string(msg.Runes)
		} else if msg.Type == tea.KeySpace {
			m.inputBuf += " "
		}
	}
	return m, nil
}

func (m Model) handleProfileKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	count := len(m.profiles) + 1 // +1 for "New Profile"
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
	case tea.KeyDown:
		if m.selectedIdx < count-1 {
			m.selectedIdx++
		}
	case tea.KeyEnter:
		if m.selectedIdx == len(m.profiles) {
			// "New Profile" selected
			m.state = StateNameInput
			m.nameInput = ""
		} else {
			// Existing profile selected
			m.player = m.profiles[m.selectedIdx]
			return m.startGame()
		}
	}
	return m, nil
}

func (m Model) handleTierKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	tiers := []string{"beginner", "explorer", "master"}
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
	case tea.KeyDown:
		if m.selectedIdx < len(tiers)-1 {
			m.selectedIdx++
		}
	case tea.KeyEnter:
		tier := tiers[m.selectedIdx]
		player, err := m.db.CreatePlayer(m.nameInput, tier)
		if err != nil {
			return m, nil
		}
		m.player = player
		return m.startGame()
	}
	return m, nil
}

func (m Model) handleNameKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		if len(strings.TrimSpace(m.nameInput)) > 0 {
			m.state = StateTierSelect
			m.selectedIdx = 0
		}
	case tea.KeyBackspace:
		if len(m.nameInput) > 0 {
			m.nameInput = m.nameInput[:len(m.nameInput)-1]
		}
	default:
		if msg.Type == tea.KeyRunes {
			m.nameInput += string(msg.Runes)
		} else if msg.Type == tea.KeySpace {
			m.nameInput += " "
		}
	}
	return m, nil
}

// startGame initializes the game with the selected player.
func (m Model) startGame() (Model, tea.Cmd) {
	w, err := world.LoadWorld("skull_island")
	if err != nil {
		return m, tea.Quit
	}
	mission := w.Missions[0]

	fs := shell.NewFS()
	paths := make([]string, 0, len(mission.Filesystem))
	for p := range mission.Filesystem {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		entry := mission.Filesystem[p]
		if entry.Type == "dir" {
			fs.Mkdir(p, entry.Hidden)
		} else {
			fs.WriteFile(p, entry.Content, entry.Hidden)
		}
	}

	ex := shell.NewExecutor(fs)
	RegisterCommands(ex, m.player.Tier)
	runner := world.NewMissionRunner(mission)

	m.fs = fs
	m.executor = ex
	m.runner = runner
	m.cwd = "/island"
	m.clueText = mission.StartingClue
	m.storyText = "Welcome to Skull Island, young pirate! Arr!"
	m.state = StateGame
	m.outputLines = nil
	m.failCount = 0
	return m, nil
}

func (m Model) submitCommand() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.inputBuf)
	m.inputBuf = ""

	if input == "" {
		return m, nil
	}

	// Echo the command
	m.outputLines = append(m.outputLines, PromptStyle.Render("$ ")+input)

	// Execute
	result := m.executor.Execute(input, m.cwd)

	if result.Clear {
		m.outputLines = nil
		return m, nil
	}

	if result.Error != "" {
		m.outputLines = append(m.outputLines, ErrorStyle.Render(result.Error))
		m.failCount++
	} else {
		if result.Output != "" {
			for _, line := range strings.Split(result.Output, "\n") {
				m.outputLines = append(m.outputLines, OutputStyle.Render(line))
			}
		}
		if result.NewCWD != "" {
			m.cwd = result.NewCWD
		}
		m.failCount = 0
	}

	// Check mission objective
	if m.runner != nil && result.Event != nil {
		if m.runner.HandleEvent(result.Event) {
			if m.runner.IsComplete() {
				m.outputLines = append(m.outputLines, SuccessStyle.Render(m.runner.Mission().SuccessMessage))
				m.storyText = m.runner.Mission().Treasure
				m.clueText = "Mission complete! You earned 3 stars!"
			} else {
				// Update clue to next objective hint
				obj := m.runner.CurrentObjective()
				if obj != nil {
					m.clueText = "Objective " + string(rune('0'+m.runner.CurrentObjectiveIndex())) + " complete! Keep going..."
				}
			}
		}
	}

	// Hint after 3 failures
	if m.failCount >= 3 && m.runner != nil {
		obj := m.runner.CurrentObjective()
		if obj != nil {
			m.storyText = "Psst! Try using '" + obj.Command + "' to progress..."
			m.failCount = 0
		}
	}

	// Trim old lines
	if len(m.outputLines) > m.maxLines {
		m.outputLines = m.outputLines[len(m.outputLines)-m.maxLines:]
	}

	return m, nil
}

// RegisterCommands registers all commands with the executor based on player tier.
func RegisterCommands(ex *shell.Executor, tier string) {
	// Beginner always
	ex.Register(commands.NewPwd())
	ex.Register(commands.NewLs())
	ex.Register(commands.NewCd())
	ex.Register(commands.NewCat())
	ex.Register(commands.NewEcho())
	ex.Register(commands.NewClear())
	ex.Register(commands.NewHelp(tier))

	if tier == "explorer" || tier == "master" {
		ex.Register(commands.NewMkdir())
		ex.Register(commands.NewTouch())
		ex.Register(commands.NewCp())
		ex.Register(commands.NewMv())
		ex.Register(commands.NewRm())
		ex.Register(commands.NewFind())
	}

	if tier == "master" {
		ex.Register(commands.NewGrep())
		ex.Register(commands.NewChmod())
		ex.Register(commands.NewMan())
		ex.Register(commands.NewHistory(func() []string { return ex.History }))
	}
}
